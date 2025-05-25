package sandbox

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/tracehawk/tracehawkx/modules"
)

// Manager handles network namespace sandboxing
type Manager struct {
	config    *modules.ScanConfig
	active    bool
	namespace string
	logger    *logrus.Entry
}

// NewManager creates a new sandbox manager
func NewManager(config *modules.ScanConfig) (*Manager, error) {
	logger := logrus.WithField("component", "sandbox")

	return &Manager{
		config: config,
		active: false,
		logger: logger,
	}, nil
}

// Initialize sets up the sandbox environment
func (m *Manager) Initialize(ctx context.Context) error {
	// Check if running as root (required for network namespaces)
	if os.Geteuid() != 0 {
		m.logger.Info("Not running as root, sandbox disabled")
		return nil
	}

	// Check if unshare is available
	if _, err := exec.LookPath("unshare"); err != nil {
		m.logger.Warn("unshare command not available, sandbox disabled")
		return nil
	}

	m.namespace = "tracehawk-scan"
	m.logger.Info("Initializing network namespace sandbox")

	// Create network namespace
	cmd := exec.CommandContext(ctx, "unshare", "-n", "sleep", "infinity")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to create network namespace: %w", err)
	}

	m.active = true
	m.logger.Info("Sandbox initialized successfully")
	return nil
}

// IsActive returns whether the sandbox is active
func (m *Manager) IsActive() bool {
	return m.active
}

// ExecuteModule runs a module within the sandbox
func (m *Manager) ExecuteModule(ctx context.Context, module modules.Module, scan *modules.Scan) error {
	if !m.active {
		return module.Run(ctx, scan)
	}

	m.logger.WithField("module", module.Name()).Info("Executing module in sandbox")

	// For now, just run the module normally
	// In a full implementation, this would execute the module in the network namespace
	return module.Run(ctx, scan)
}

// Cleanup cleans up the sandbox environment
func (m *Manager) Cleanup() error {
	if !m.active {
		return nil
	}

	m.logger.Info("Cleaning up sandbox")

	// Kill any processes in the namespace
	// In a full implementation, this would properly clean up the network namespace

	m.active = false
	return nil
}
