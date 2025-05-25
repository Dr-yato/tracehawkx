package stable

import (
	"context"
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tracehawk/tracehawkx/modules"
)

// NmapModule implements port scanning using nmap
type NmapModule struct{}

func (m *NmapModule) Name() string {
	return "nmap"
}

func (m *NmapModule) Description() string {
	return "Network port scanner and service detection"
}

func (m *NmapModule) Category() string {
	return "stable"
}

func (m *NmapModule) Author() string {
	return "TraceHawk Team"
}

func (m *NmapModule) Version() string {
	return "1.0.0"
}

func (m *NmapModule) Flags(fs *flag.FlagSet) {
	// Module-specific flags can be added here
}

func (m *NmapModule) Prerequisites() error {
	// Check if nmap is available
	if _, err := exec.LookPath("nmap"); err != nil {
		return fmt.Errorf("nmap not found in PATH")
	}
	return nil
}

func (m *NmapModule) Run(ctx context.Context, scan *modules.Scan) error {
	scan.Logger.Info("Running nmap for port scanning")

	for _, target := range scan.Targets {
		// Build nmap command
		args := []string{"-sS", "-sV", "-O", "--top-ports", "1000", "-T4", target}
		if scan.Config.Stealth {
			args = []string{"-sS", "-T2", target}
		}
		if scan.Config.Aggressive {
			args = []string{"-A", "-T5", target}
		}

		cmd := exec.CommandContext(ctx, "nmap", args...)
		output, err := cmd.Output()
		if err != nil {
			scan.Logger.Warnf("Nmap failed for %s: %v", target, err)
			continue
		}

		// Parse nmap output (simplified)
		host := m.parseNmapOutput(string(output), target)
		if host != nil {
			scan.AddHost(*host)
		}

		scan.Logger.Infof("Scanned %s with %d open ports", target, len(host.Ports))
	}

	return nil
}

func (m *NmapModule) parseNmapOutput(output, target string) *modules.HostResult {
	host := &modules.HostResult{
		IP:       target,
		Hostname: target,
		Ports:    []modules.PortResult{},
		Services: []modules.ServiceResult{},
		WebApps:  []modules.WebAppResult{},
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse port lines (simplified)
		if strings.Contains(line, "/tcp") && strings.Contains(line, "open") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				portStr := strings.Split(parts[0], "/")[0]
				if port, err := strconv.Atoi(portStr); err == nil {
					portResult := modules.PortResult{
						Port:     port,
						Protocol: "tcp",
						State:    "open",
						Service:  parts[2],
						Version:  "",
					}
					if len(parts) > 3 {
						portResult.Version = strings.Join(parts[3:], " ")
					}
					host.Ports = append(host.Ports, portResult)
				}
			}
		}
	}

	return host
}

func (m *NmapModule) Cleanup() error {
	return nil
}
