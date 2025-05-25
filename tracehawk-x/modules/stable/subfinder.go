package stable

import (
	"context"
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/tracehawk/tracehawkx/modules"
)

// SubfinderModule implements subdomain discovery using subfinder
type SubfinderModule struct{}

func (m *SubfinderModule) Name() string {
	return "subfinder"
}

func (m *SubfinderModule) Description() string {
	return "Fast passive subdomain enumeration using multiple sources"
}

func (m *SubfinderModule) Category() string {
	return "stable"
}

func (m *SubfinderModule) Author() string {
	return "TraceHawk Team"
}

func (m *SubfinderModule) Version() string {
	return "1.0.0"
}

func (m *SubfinderModule) Flags(fs *flag.FlagSet) {
	// Module-specific flags can be added here
}

func (m *SubfinderModule) Prerequisites() error {
	// Check if subfinder is available
	if _, err := exec.LookPath("subfinder"); err != nil {
		return fmt.Errorf("subfinder not found in PATH")
	}
	return nil
}

func (m *SubfinderModule) Run(ctx context.Context, scan *modules.Scan) error {
	scan.Logger.Info("Running subfinder for subdomain discovery")

	for _, target := range scan.Targets {
		// Skip if target is an IP address
		if strings.Contains(target, "/") || strings.Count(target, ".") == 3 {
			continue
		}

		// Run subfinder command
		cmd := exec.CommandContext(ctx, "subfinder", "-d", target, "-silent")
		output, err := cmd.Output()
		if err != nil {
			scan.Logger.Warnf("Subfinder failed for %s: %v", target, err)
			continue
		}

		// Parse subdomains
		subdomains := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, subdomain := range subdomains {
			if subdomain != "" {
				// Add discovered subdomain as a host
				host := modules.HostResult{
					IP:       "",
					Hostname: subdomain,
					Ports:    []modules.PortResult{},
					Services: []modules.ServiceResult{},
					WebApps:  []modules.WebAppResult{},
				}
				scan.AddHost(host)
			}
		}

		scan.Logger.Infof("Discovered %d subdomains for %s", len(subdomains), target)
	}

	return nil
}

func (m *SubfinderModule) Cleanup() error {
	return nil
}
