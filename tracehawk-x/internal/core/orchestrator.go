package core

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/tracehawk/tracehawkx/internal/report"
	"github.com/tracehawk/tracehawkx/internal/sandbox"
	"github.com/tracehawk/tracehawkx/internal/scoring"
	"github.com/tracehawk/tracehawkx/modules"
)

// Orchestrator manages the overall scan execution
type Orchestrator struct {
	config   *modules.ScanConfig
	scan     *modules.Scan
	logger   *logrus.Entry
	sandbox  *sandbox.Manager
	reporter *report.Generator
	scorer   *scoring.Engine
}

// NewOrchestrator creates a new scan orchestrator
func NewOrchestrator(config *modules.ScanConfig) (*Orchestrator, error) {
	scanID := uuid.New().String()

	// Create logger
	logger := logrus.WithFields(logrus.Fields{
		"scan_id": scanID,
		"targets": len(config.Targets),
	})

	// Initialize scan context
	scan := &modules.Scan{
		ID:      scanID,
		Targets: config.Targets,
		Config:  config,
		Results: &modules.ScanResults{
			BleedingEdge: make(map[string]interface{}),
		},
		Context: make(map[string]interface{}),
		Logger:  logger,
	}

	// Initialize sandbox manager
	sandboxMgr, err := sandbox.NewManager(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sandbox: %w", err)
	}

	// Initialize report generator
	reportGen, err := report.NewGenerator(config.ReportDir, config.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize report generator: %w", err)
	}

	// Initialize scoring engine
	scoringEngine := scoring.NewEngine()

	return &Orchestrator{
		config:   config,
		scan:     scan,
		logger:   logger,
		sandbox:  sandboxMgr,
		reporter: reportGen,
		scorer:   scoringEngine,
	}, nil
}

// Execute runs the complete scan workflow
func (o *Orchestrator) Execute(ctx context.Context) error {
	startTime := time.Now()
	o.logger.Info("🦅 Starting TraceHawk X scan")

	defer func() {
		duration := time.Since(startTime)
		o.scan.Results.Summary.ScanDuration = duration.String()
		o.logger.WithField("duration", duration).Info("Scan completed")
	}()

	// Validate targets
	if err := o.validateTargets(); err != nil {
		return fmt.Errorf("target validation failed: %w", err)
	}

	// Initialize sandbox if not running as root
	if err := o.sandbox.Initialize(ctx); err != nil {
		o.logger.Warnf("Sandbox initialization failed, continuing without sandbox: %v", err)
	}
	defer o.sandbox.Cleanup()

	// Execute scan phases
	phases := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"Discovery", o.runDiscoveryPhase},
		{"Scanning", o.runScanningPhase},
		{"Vulnerability Assessment", o.runVulnPhase},
		{"Bleeding Edge", o.runBleedingEdgePhase},
		{"Scoring", o.runScoringPhase},
		{"Reporting", o.runReportingPhase},
	}

	for _, phase := range phases {
		o.logger.Infof("📍 Phase: %s", phase.name)

		phaseStart := time.Now()
		if err := phase.fn(ctx); err != nil {
			return fmt.Errorf("phase %s failed: %w", phase.name, err)
		}

		o.logger.WithField("duration", time.Since(phaseStart)).
			Infof("✅ Phase %s completed", phase.name)

		// Check for cancellation between phases
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	o.logger.Info("🎯 Scan execution completed successfully")
	return nil
}

// runDiscoveryPhase executes asset discovery modules
func (o *Orchestrator) runDiscoveryPhase(ctx context.Context) error {
	discoveryModules := []string{
		"subfinder",
		"amass",
		"dnsx",
		"asnmap",
	}

	return o.executeModules(ctx, discoveryModules, "stable")
}

// runScanningPhase executes port and service scanning
func (o *Orchestrator) runScanningPhase(ctx context.Context) error {
	scanningModules := []string{
		"naabu",
		"nmap",
		"httpx",
	}

	return o.executeModules(ctx, scanningModules, "stable")
}

// runVulnPhase executes vulnerability scanning
func (o *Orchestrator) runVulnPhase(ctx context.Context) error {
	vulnModules := []string{
		"nuclei",
		"katana",
		"ffuf",
	}

	return o.executeModules(ctx, vulnModules, "stable")
}

// runBleedingEdgePhase executes bleeding-edge modules if enabled
func (o *Orchestrator) runBleedingEdgePhase(ctx context.Context) error {
	if !o.config.BleedingEdge {
		o.logger.Info("Skipping bleeding-edge modules (not enabled)")
		return nil
	}

	bleedingModules := []string{
		"llm-fuzzer",
		"auto-patch",
		"shadow-clone",
		"dep-drift",
		"timing-map",
		"blue-team",
	}

	return o.executeModules(ctx, bleedingModules, "bleeding-edge")
}

// runScoringPhase calculates risk scores
func (o *Orchestrator) runScoringPhase(ctx context.Context) error {
	o.logger.Info("📊 Calculating risk scores")

	// Score all vulnerabilities
	for i := range o.scan.Results.Vulnerabilities {
		vuln := &o.scan.Results.Vulnerabilities[i]
		score := o.scorer.CalculateRiskScore(vuln, o.scan)
		vuln.RiskScore = score

		// Mark high-risk vulnerabilities
		if score >= 70.0 {
			o.scan.Results.Summary.HighRiskVulns++
		}
	}

	// Calculate overall scan risk score
	overallScore := o.scorer.CalculateOverallRiskScore(o.scan.Results.Vulnerabilities)
	o.scan.Results.Summary.RiskScore = overallScore

	// Update summary statistics
	o.updateScanSummary()

	return nil
}

// runReportingPhase generates scan reports
func (o *Orchestrator) runReportingPhase(ctx context.Context) error {
	o.logger.Info("📝 Generating reports")

	return o.reporter.Generate(ctx, o.scan)
}

// executeModules runs a set of modules concurrently with proper error handling
func (o *Orchestrator) executeModules(ctx context.Context, moduleNames []string, category string) error {
	// Get available modules for this category
	availableModules := modules.GetModulesByCategory(category)
	moduleMap := make(map[string]modules.Module)
	for _, mod := range availableModules {
		moduleMap[mod.Name()] = mod
	}

	// Filter requested modules to only available ones
	var activeModules []modules.Module
	for _, name := range moduleNames {
		if mod, exists := moduleMap[name]; exists {
			// Check prerequisites
			if err := mod.Prerequisites(); err != nil {
				o.logger.Warnf("Module %s prerequisites not met, skipping: %v", name, err)
				continue
			}
			activeModules = append(activeModules, mod)
		}
	}

	if len(activeModules) == 0 {
		o.logger.Warnf("No modules available for category %s", category)
		return nil
	}

	// Execute modules concurrently with limited parallelism
	g, gCtx := errgroup.WithContext(ctx)

	// Limit concurrency to prevent resource exhaustion
	sem := make(chan struct{}, 3)

	for _, mod := range activeModules {
		mod := mod // Capture loop variable

		g.Go(func() error {
			// Acquire semaphore
			select {
			case sem <- struct{}{}:
			case <-gCtx.Done():
				return gCtx.Err()
			}
			defer func() { <-sem }()

			return o.executeModule(gCtx, mod)
		})
	}

	return g.Wait()
}

// executeModule runs a single module with proper logging and error handling
func (o *Orchestrator) executeModule(ctx context.Context, mod modules.Module) error {
	logger := o.logger.WithField("module", mod.Name())
	logger.Info("🔧 Executing module")

	start := time.Now()
	defer func() {
		duration := time.Since(start)
		logger.WithField("duration", duration).Info("Module execution completed")
	}()

	// Execute module in sandbox if available
	var err error
	if o.sandbox.IsActive() {
		err = o.sandbox.ExecuteModule(ctx, mod, o.scan)
	} else {
		err = mod.Run(ctx, o.scan)
	}

	if err != nil {
		logger.Errorf("Module execution failed: %v", err)
		// Don't fail entire scan for single module failure
		return nil
	}

	// Track executed modules
	o.scan.Results.Summary.ModulesExecuted = append(
		o.scan.Results.Summary.ModulesExecuted,
		mod.Name(),
	)

	return nil
}

// validateTargets validates the target specification
func (o *Orchestrator) validateTargets() error {
	if len(o.scan.Targets) == 0 {
		return fmt.Errorf("no targets specified")
	}

	return modules.ValidateTargets(o.scan.Targets)
}

// updateScanSummary updates the scan summary statistics
func (o *Orchestrator) updateScanSummary() {
	summary := &o.scan.Results.Summary

	// Count hosts
	summary.TotalHosts = len(o.scan.Results.Hosts)
	summary.AliveHosts = 0
	for _, host := range o.scan.Results.Hosts {
		if len(host.Ports) > 0 || len(host.Services) > 0 {
			summary.AliveHosts++
		}
	}

	// Count ports
	summary.TotalPorts = 0
	summary.OpenPorts = 0
	for _, host := range o.scan.Results.Hosts {
		summary.TotalPorts += len(host.Ports)
		for _, port := range host.Ports {
			if port.State == "open" {
				summary.OpenPorts++
			}
		}
	}

	// Count vulnerabilities by severity
	summary.TotalVulns = len(o.scan.Results.Vulnerabilities)
	summary.CriticalVulns = 0
	summary.HighVulns = 0
	summary.MediumVulns = 0
	summary.LowVulns = 0
	summary.ExploitableVulns = 0

	for _, vuln := range o.scan.Results.Vulnerabilities {
		switch vuln.Severity {
		case "critical":
			summary.CriticalVulns++
		case "high":
			summary.HighVulns++
		case "medium":
			summary.MediumVulns++
		case "low":
			summary.LowVulns++
		}

		if vuln.Exploitable {
			summary.ExploitableVulns++
		}
	}
}
