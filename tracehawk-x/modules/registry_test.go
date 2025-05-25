package modules

import (
	"context"
	"flag"
	"testing"
)

// MockModule for testing
type MockModule struct {
	name        string
	description string
	category    string
}

func (m *MockModule) Name() string                              { return m.name }
func (m *MockModule) Description() string                       { return m.description }
func (m *MockModule) Category() string                          { return m.category }
func (m *MockModule) Author() string                            { return "Test Author" }
func (m *MockModule) Version() string                           { return "1.0.0" }
func (m *MockModule) Flags(fs *flag.FlagSet)                    {}
func (m *MockModule) Prerequisites() error                      { return nil }
func (m *MockModule) Run(ctx context.Context, scan *Scan) error { return nil }
func (m *MockModule) Cleanup() error                            { return nil }

func TestModuleRegistry(t *testing.T) {
	// Create a mock module
	mockModule := &MockModule{
		name:        "test-module",
		description: "Test module for unit testing",
		category:    "stable",
	}

	// Register the module
	Register(mockModule)

	// Test GetModule
	retrieved, exists := GetModule("test-module")
	if !exists {
		t.Error("Expected module to exist after registration")
	}

	if retrieved.Name() != "test-module" {
		t.Errorf("Expected module name 'test-module', got '%s'", retrieved.Name())
	}

	// Test GetModulesByCategory
	stableModules := GetModulesByCategory("stable")
	found := false
	for _, mod := range stableModules {
		if mod.Name() == "test-module" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find test-module in stable category")
	}
}

func TestScanHelperMethods(t *testing.T) {
	scan := &Scan{
		Results: &ScanResults{},
		Context: make(map[string]interface{}),
	}

	// Test AddHost
	host := HostResult{
		IP:       "192.168.1.1",
		Hostname: "test.example.com",
	}
	scan.AddHost(host)

	if len(scan.Results.Hosts) != 1 {
		t.Errorf("Expected 1 host, got %d", len(scan.Results.Hosts))
	}

	// Test AddVulnerability
	vuln := Vulnerability{
		ID:       "TEST-001",
		Name:     "Test Vulnerability",
		Severity: "high",
	}
	scan.AddVulnerability(vuln)

	if len(scan.Results.Vulnerabilities) != 1 {
		t.Errorf("Expected 1 vulnerability, got %d", len(scan.Results.Vulnerabilities))
	}

	// Test Context methods
	scan.SetContext("test_key", "test_value")
	value, exists := scan.GetContext("test_key")
	if !exists {
		t.Error("Expected context key to exist")
	}
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%v'", value)
	}
}

func TestUtilityFunctions(t *testing.T) {
	// Test EscapeShellArg
	input := `test"arg\with'quotes`
	escaped := EscapeShellArg(input)
	if escaped == input {
		t.Error("Expected input to be escaped")
	}

	// Test ValidateTargets (placeholder)
	targets := []string{"example.com", "192.168.1.1"}
	err := ValidateTargets(targets)
	if err != nil {
		t.Errorf("ValidateTargets failed: %v", err)
	}
}
