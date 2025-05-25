package scoring

import (
	"math"

	"github.com/tracehawk/tracehawkx/modules"
)

// Engine calculates risk scores for vulnerabilities
type Engine struct {
	// Configuration for scoring weights
	severityWeights map[string]float64
	exploitWeight   float64
	driftWeight     float64
	llmWeight       float64
}

// NewEngine creates a new scoring engine
func NewEngine() *Engine {
	return &Engine{
		severityWeights: map[string]float64{
			"critical": 10.0,
			"high":     7.5,
			"medium":   5.0,
			"low":      2.5,
			"info":     1.0,
		},
		exploitWeight: 1.5,
		driftWeight:   0.2,
		llmWeight:     0.2,
	}
}

// CalculateRiskScore calculates the risk score for a vulnerability
// Formula: RiskScore = (BaseSeverity × 10) · ExploitabilityCoeff · (1 + SupplyChainDrift) · LLMConfidence
func (e *Engine) CalculateRiskScore(vuln *modules.Vulnerability, scan *modules.Scan) float64 {
	// Base severity score
	baseSeverity := e.severityWeights[vuln.Severity]
	if baseSeverity == 0 {
		baseSeverity = 1.0 // Default for unknown severity
	}

	// Exploitability coefficient
	exploitCoeff := 1.0
	if vuln.Exploitable {
		exploitCoeff = e.exploitWeight
	}

	// Supply chain drift factor
	driftFactor := 1.0
	if driftData, exists := scan.GetContext("supply_chain_drift"); exists {
		if drift, ok := driftData.(float64); ok {
			driftFactor = 1.0 + (drift * e.driftWeight)
		}
	}

	// LLM confidence factor
	llmConfidence := 1.0
	if confidence, exists := vuln.Metadata["llm_confidence"]; exists {
		if conf, ok := confidence.(float64); ok {
			llmConfidence = 0.8 + (conf * e.llmWeight)
		}
	}

	// Calculate final score
	score := baseSeverity * exploitCoeff * driftFactor * llmConfidence

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return math.Round(score*100) / 100
}

// CalculateOverallRiskScore calculates the overall risk score for all vulnerabilities
func (e *Engine) CalculateOverallRiskScore(vulns []modules.Vulnerability) float64 {
	if len(vulns) == 0 {
		return 0.0
	}

	var totalScore float64
	var criticalCount, highCount int

	for _, vuln := range vulns {
		totalScore += vuln.RiskScore

		switch vuln.Severity {
		case "critical":
			criticalCount++
		case "high":
			highCount++
		}
	}

	// Average score with bonus for critical/high vulnerabilities
	avgScore := totalScore / float64(len(vulns))

	// Bonus for having critical vulnerabilities
	criticalBonus := float64(criticalCount) * 5.0
	highBonus := float64(highCount) * 2.0

	overallScore := avgScore + criticalBonus + highBonus

	// Cap at 100
	if overallScore > 100 {
		overallScore = 100
	}

	return math.Round(overallScore*100) / 100
}
