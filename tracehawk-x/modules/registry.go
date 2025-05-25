package modules

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"strings"
	"sync"

	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

// Module represents a scanning module plugin interface
type Module interface {
	// Name returns the module name
	Name() string

	// Description returns a short description
	Description() string

	// Category returns the module category (stable, bleeding-edge)
	Category() string

	// Author returns the module author
	Author() string

	// Version returns the module version
	Version() string

	// Flags registers module-specific command line flags
	Flags(fs *flag.FlagSet)

	// Prerequisites checks if module dependencies are met
	Prerequisites() error

	// Run executes the module logic
	Run(ctx context.Context, scan *Scan) error

	// Cleanup performs any necessary cleanup
	Cleanup() error
}

// Scan represents the current scan context and shared data
type Scan struct {
	ID      string                 `json:"id"`
	Targets []string               `json:"targets"`
	Config  *ScanConfig            `json:"config"`
	Results *ScanResults           `json:"results"`
	Context map[string]interface{} `json:"context"`
	Logger  *logrus.Entry          `json:"-"`

	mu sync.RWMutex
}

// ScanConfig holds scan configuration
type ScanConfig struct {
	Targets      []string `json:"targets"`
	Output       string   `json:"output"`
	ReportDir    string   `json:"report_dir"`
	Threads      int      `json:"threads"`
	RateLimit    int      `json:"rate_limit"`
	Deep         bool     `json:"deep"`
	BleedingEdge bool     `json:"bleeding_edge"`
	Stealth      bool     `json:"stealth"`
	Aggressive   bool     `json:"aggressive"`
	NoThrottle   bool     `json:"no_throttle"`
	Exclude      []string `json:"exclude"`

	// Module-specific config
	LLMModel      string  `json:"llm_model"`
	Temperature   float64 `json:"temperature"`
	GeneratePatch bool    `json:"generate_patch"`
	ShadowClone   string  `json:"shadow_clone"`
	DepDrift      bool    `json:"dep_drift"`
	TimingMap     bool    `json:"timing_map"`
	BlueTeam      bool    `json:"blue_team"`
}

// ScanResults holds aggregated scan results
type ScanResults struct {
	Hosts           []HostResult           `json:"hosts"`
	Vulnerabilities []Vulnerability        `json:"vulnerabilities"`
	Patches         []PatchRecommendation  `json:"patches"`
	BleedingEdge    map[string]interface{} `json:"bleeding_edge"`
	Summary         ScanSummary            `json:"summary"`

	mu sync.RWMutex
}

// HostResult represents a scanned host
type HostResult struct {
	IP       string          `json:"ip"`
	Hostname string          `json:"hostname"`
	Ports    []PortResult    `json:"ports"`
	Services []ServiceResult `json:"services"`
	WebApps  []WebAppResult  `json:"webapps"`
}

// PortResult represents a discovered port
type PortResult struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
	Service  string `json:"service"`
	Version  string `json:"version"`
}

// ServiceResult represents a discovered service
type ServiceResult struct {
	Name     string                 `json:"name"`
	Version  string                 `json:"version"`
	Banner   string                 `json:"banner"`
	Metadata map[string]interface{} `json:"metadata"`
}

// WebAppResult represents a discovered web application
type WebAppResult struct {
	URL        string                 `json:"url"`
	StatusCode int                    `json:"status_code"`
	Title      string                 `json:"title"`
	Technology []string               `json:"technology"`
	Headers    map[string]string      `json:"headers"`
	Endpoints  []string               `json:"endpoints"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// Vulnerability represents a discovered vulnerability
type Vulnerability struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"`
	CVSS        float64                `json:"cvss"`
	CVE         string                 `json:"cve"`
	Host        string                 `json:"host"`
	Port        int                    `json:"port"`
	Service     string                 `json:"service"`
	URL         string                 `json:"url"`
	Evidence    string                 `json:"evidence"`
	PoC         string                 `json:"poc"`
	References  []string               `json:"references"`
	RiskScore   float64                `json:"risk_score"`
	Exploitable bool                   `json:"exploitable"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// PatchRecommendation represents an auto-generated patch
type PatchRecommendation struct {
	VulnID      string  `json:"vuln_id"`
	Type        string  `json:"type"` // "code", "config", "waf"
	Language    string  `json:"language"`
	Framework   string  `json:"framework"`
	Diff        string  `json:"diff"`
	WAFRule     string  `json:"waf_rule"`
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"`
}

// ScanSummary provides scan statistics
type ScanSummary struct {
	TotalHosts       int      `json:"total_hosts"`
	AliveHosts       int      `json:"alive_hosts"`
	TotalPorts       int      `json:"total_ports"`
	OpenPorts        int      `json:"open_ports"`
	TotalVulns       int      `json:"total_vulns"`
	CriticalVulns    int      `json:"critical_vulns"`
	HighVulns        int      `json:"high_vulns"`
	MediumVulns      int      `json:"medium_vulns"`
	LowVulns         int      `json:"low_vulns"`
	ExploitableVulns int      `json:"exploitable_vulns"`
	HighRiskVulns    int      `json:"high_risk_vulns"`
	ScanDuration     string   `json:"scan_duration"`
	RiskScore        float64  `json:"risk_score"`
	ModulesExecuted  []string `json:"modules_executed"`
}

// Registry holds all registered modules
type Registry struct {
	modules map[string]Module
	mu      sync.RWMutex
}

var (
	registry = &Registry{
		modules: make(map[string]Module),
	}
)

// Register registers a new module
func Register(module Module) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	name := module.Name()
	if name == "" {
		logrus.Fatalf("Module name cannot be empty")
	}

	if _, exists := registry.modules[name]; exists {
		logrus.Fatalf("Module %s already registered", name)
	}

	registry.modules[name] = module
	logrus.Debugf("Registered module: %s", name)
}

// GetModule retrieves a module by name
func GetModule(name string) (Module, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	module, exists := registry.modules[name]
	return module, exists
}

// GetModules returns all registered modules
func GetModules() map[string]Module {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	modules := make(map[string]Module)
	for name, module := range registry.modules {
		modules[name] = module
	}
	return modules
}

// GetModulesByCategory returns modules filtered by category
func GetModulesByCategory(category string) []Module {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	var filtered []Module
	for _, module := range registry.modules {
		if module.Category() == category {
			filtered = append(filtered, module)
		}
	}
	return filtered
}

// ListModules displays all registered modules in a table format
func ListModules() {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	if len(registry.modules) == 0 {
		fmt.Println(color.YellowString("No modules registered"))
		return
	}

	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Category", "Version", "Author", "Description"})
	table.SetBorder(false)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
	)

	// Sort modules by category and name
	var names []string
	for name := range registry.modules {
		names = append(names, name)
	}
	sort.Strings(names)

	// Group by category
	categories := make(map[string][]string)
	for _, name := range names {
		module := registry.modules[name]
		category := module.Category()
		categories[category] = append(categories[category], name)
	}

	// Add rows grouped by category
	for _, category := range []string{"stable", "bleeding-edge"} {
		if moduleNames, exists := categories[category]; exists {
			// Add category separator
			if category == "bleeding-edge" {
				table.Append([]string{"", "", "", "", ""})
			}

			for _, name := range moduleNames {
				module := registry.modules[name]

				// Color code by category
				var nameColor, categoryColor string
				if category == "stable" {
					nameColor = color.GreenString(name)
					categoryColor = color.GreenString(category)
				} else {
					nameColor = color.MagentaString(name)
					categoryColor = color.MagentaString(category)
				}

				// Truncate description if too long
				desc := module.Description()
				if len(desc) > 60 {
					desc = desc[:57] + "..."
				}

				table.Append([]string{
					nameColor,
					categoryColor,
					module.Version(),
					module.Author(),
					desc,
				})
			}
		}
	}

	fmt.Println(color.HiCyanString("📋 Available Modules"))
	fmt.Println()
	table.Render()
	fmt.Println()

	// Show usage hints
	fmt.Println(color.HiYellowString("💡 Usage:"))
	fmt.Println("  • Use --bleeding-edge to enable bleeding-edge modules")
	fmt.Println("  • Use --deep for thorough scanning with all stable modules")
	fmt.Println("  • See docs/MODULES.md for detailed module documentation")
}

// Helper methods for Scan
func (s *Scan) AddHost(host HostResult) {
	s.Results.mu.Lock()
	defer s.Results.mu.Unlock()
	s.Results.Hosts = append(s.Results.Hosts, host)
}

func (s *Scan) AddVulnerability(vuln Vulnerability) {
	s.Results.mu.Lock()
	defer s.Results.mu.Unlock()
	s.Results.Vulnerabilities = append(s.Results.Vulnerabilities, vuln)
}

func (s *Scan) AddPatch(patch PatchRecommendation) {
	s.Results.mu.Lock()
	defer s.Results.mu.Unlock()
	s.Results.Patches = append(s.Results.Patches, patch)
}

func (s *Scan) SetBleedingEdgeResult(moduleID string, result interface{}) {
	s.Results.mu.Lock()
	defer s.Results.mu.Unlock()
	if s.Results.BleedingEdge == nil {
		s.Results.BleedingEdge = make(map[string]interface{})
	}
	s.Results.BleedingEdge[moduleID] = result
}

func (s *Scan) GetContext(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, exists := s.Context[key]
	return val, exists
}

func (s *Scan) SetContext(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Context == nil {
		s.Context = make(map[string]interface{})
	}
	s.Context[key] = value
}

// Utility functions for modules
func IsCommandAvailable(command string) bool {
	// This will be implemented to check if external tools are available
	return true // Placeholder
}

func ValidateTargets(targets []string) error {
	// This will be implemented to validate target formats
	return nil // Placeholder
}

func EscapeShellArg(arg string) string {
	// Escape shell arguments to prevent injection
	return strings.ReplaceAll(strings.ReplaceAll(arg, `\`, `\\`), `"`, `\"`)
}
