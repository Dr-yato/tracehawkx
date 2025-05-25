# 🦅 TraceHawk X

**Autonomous Terminal-Native Recon & Exploitation Agent**

TraceHawk X is a comprehensive penetration testing toolkit that performs end-to-end security assessments from enumeration through exploitation to auto-patch suggestions, without relying on any cloud APIs.

🌐 **Visit [hunter3.ninja](https://hunter3.ninja) for more security tools and resources**

## ✨ Features

- **🔒 Fully Offline**: No external APIs required, works completely offline
- **🧩 Modular Architecture**: Plugin-based system with stable and bleeding-edge modules
- **📊 Beautiful Reports**: Generates Markdown and HTML reports with executive summaries
- **🛡️ Network Sandboxing**: Runs scanners in isolated network namespaces
- **🤖 AI-Powered**: Local LLM integration for intelligent fuzzing and patch generation
- **⚡ High Performance**: Concurrent execution with intelligent throttling
- **🎯 Risk Scoring**: Advanced vulnerability prioritization with CVSS integration

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/hunter3ninja/tracehawkx.git
cd tracehawkx

# Build the binary
make build

# Install system-wide (optional)
make install
```

### Basic Usage

```bash
# Initialize configuration
tracehawkx init

# Basic scan
tracehawkx scan --target example.com

# Deep scan with bleeding-edge modules
tracehawkx scan --target example.com --deep --bleeding-edge

# Scan with custom output
tracehawkx scan --target example.com --output results/ --report reports/

# List available modules
tracehawkx list-modules

# Launch interactive TUI
tracehawkx tui
```

## 📋 Module Categories

### Stable Modules
- **subfinder** - Fast passive subdomain enumeration
- **amass** - In-depth DNS enumeration and network mapping
- **dnsx** - Fast DNS toolkit for resolution and discovery
- **asnmap** - ASN to CIDR mapping for network reconnaissance
- **naabu** - Fast port scanner written in Go
- **nmap** - Network exploration and security auditing
- **httpx** - Fast HTTP probing and web application discovery
- **nuclei** - Vulnerability scanner based on templates
- **katana** - Web crawling and spidering
- **ffuf** - Fast web fuzzer for content discovery

### Bleeding-Edge Modules
- **llm-fuzzer** - LLM-guided semantic fuzzer using local Llama models
- **auto-patch** - Generates code patches and WAF rules for vulnerabilities
- **shadow-clone** - Mirrors production traffic to test auth-bound flaws
- **dep-drift** - Supply chain dependency drift scanner
- **timing-map** - Side-channel timing analysis using eBPF
- **blue-team** - Generates blue team replay kits with synthetic logs

## 🔧 Configuration

TraceHawk X uses a YAML configuration file located at `~/.tracehawkx.yaml`:

```yaml
# TraceHawk X Configuration
output: results
report_dir: reports
threads: 50
rate_limit: 150
deep: false
bleeding_edge: false
stealth: false
aggressive: false
no_throttle: false

# Module-specific settings
temperature: 0.7
generate_patch: false
dep_drift: false
timing_map: false
blue_team: false

# Exclude patterns (optional)
exclude: []
```

## 📊 Risk Scoring

TraceHawk X uses an advanced risk scoring algorithm:

```
RiskScore = (BaseSeverity × 10) · ExploitabilityCoeff · (1 + SupplyChainDrift) · LLMConfidence
```

- **BaseSeverity**: CVSS-based severity weighting
- **ExploitabilityCoeff**: Bonus for confirmed exploitable vulnerabilities
- **SupplyChainDrift**: Supply chain risk distance factor
- **LLMConfidence**: AI confidence in exploit chain viability

High-risk vulnerabilities (≥70) are prioritized in the "Immediate Threats" section.

## 🐳 Docker Usage

```bash
# Build Docker image
make docker

# Run in container
docker run --rm -v $(pwd)/results:/app/results tracehawkx:latest scan --target example.com
```

## 🛠️ Development

### Prerequisites

- Go 1.22+
- Make
- Docker (optional)

### Building from Source

```bash
# Download dependencies
make deps

# Build binary
make build

# Run tests
make test

# Run linter
make lint

# Build for all platforms
make release
```

### Adding New Modules

1. Create a new module file in `modules/stable/` or `modules/bleeding-edge/`
2. Implement the `Module` interface:

```go
type Module interface {
    Name() string
    Description() string
    Category() string
    Author() string
    Version() string
    Flags(fs *flag.FlagSet)
    Prerequisites() error
    Run(ctx context.Context, scan *Scan) error
    Cleanup() error
}
```

3. Register the module in the appropriate `init.go` file
4. Rebuild and test

## 📖 Documentation

- [Usage Guide](docs/USAGE.md)
- [Module Development](docs/MODULES.md)
- [Development Guide](docs/DEVELOPMENT.md)
- [API Reference](docs/API.md)

## 🔒 Security

TraceHawk X is designed with security in mind:

- **Network Isolation**: Modules run in separate network namespaces when possible
- **Input Validation**: All user inputs are validated and sanitized
- **Rate Limiting**: Built-in throttling to prevent overwhelming targets
- **Privilege Separation**: Runs with minimal required privileges

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

TraceHawk X is released under the MIT License. See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

TraceHawk X builds upon the excellent work of:

- [ProjectDiscovery](https://projectdiscovery.io/) - For nuclei, subfinder, httpx, and other tools
- [OWASP](https://owasp.org/) - For security testing methodologies
- [The Go Team](https://golang.org/) - For the excellent Go programming language

## ⚠️ Disclaimer

TraceHawk X is intended for authorized security testing only. Users are responsible for complying with applicable laws and regulations. The authors are not responsible for any misuse of this tool.

---

**Made with ❤️ by the [Hunter3.Ninja](https://hunter3.ninja) Team** 🥷 

🦅 TraceHawk X - Autonomous Terminal-Native Recon & Exploitation Agent | Advanced penetration testing toolkit with AI-powered modules 

penetration-testing
security-tools
reconnaissance
vulnerability-scanner
golang
cybersecurity
red-team
exploit-development
llm-integration
automated-testing 
