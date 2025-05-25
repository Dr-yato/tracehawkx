# TraceHawk X Usage Guide

## Table of Contents

- [Installation](#installation)
- [Basic Commands](#basic-commands)
- [Scan Modes](#scan-modes)
- [Configuration](#configuration)
- [Output Formats](#output-formats)
- [Module Management](#module-management)
- [Advanced Usage](#advanced-usage)
- [Examples](#examples)

## Installation

### From Source

```bash
git clone https://github.com/tracehawk/tracehawkx.git
cd tracehawkx
make build
sudo make install
```

### Using Docker

```bash
docker pull tracehawkx:latest
```

### Binary Releases

Download pre-compiled binaries from the [releases page](https://github.com/tracehawk/tracehawkx/releases).

## Basic Commands

### Initialize Configuration

```bash
tracehawkx init
```

This creates a default configuration file at `~/.tracehawkx.yaml` and sets up template directories.

### Basic Scan

```bash
tracehawkx scan --target example.com
```

### List Available Modules

```bash
tracehawkx list-modules
```

### Show Version

```bash
tracehawkx version
```

### Launch TUI

```bash
tracehawkx tui
```

## Scan Modes

### Standard Scan

```bash
tracehawkx scan --target example.com
```

Runs basic reconnaissance and vulnerability scanning using stable modules.

### Deep Scan

```bash
tracehawkx scan --target example.com --deep
```

Enables thorough scanning with all stable modules and extended enumeration.

### Bleeding-Edge Scan

```bash
tracehawkx scan --target example.com --bleeding-edge
```

Activates experimental modules including LLM fuzzing and advanced analysis.

### Stealth Scan

```bash
tracehawkx scan --target example.com --stealth
```

Uses slower, less detectable scanning techniques.

### Aggressive Scan

```bash
tracehawkx scan --target example.com --aggressive
```

Fast, comprehensive scanning that may be easily detected.

## Configuration

### Global Flags

- `--config, -c`: Specify custom config file
- `--verbose, -v`: Enable verbose output
- `--debug, -d`: Enable debug output
- `--no-color`: Disable colored output

### Scan Flags

#### Target Specification
- `--target, -t`: Target hosts, domains, or CIDR ranges
- `--target-file, -T`: File containing targets (one per line)
- `--exclude`: Hosts/ranges to exclude from scan

#### Scan Modes
- `--deep`: Enable deep scanning
- `--bleeding-edge`: Enable bleeding-edge modules
- `--stealth`: Use stealth scanning techniques
- `--aggressive`: Aggressive scanning

#### Output Options
- `--output, -o`: Output directory (default: results)
- `--report, -r`: Report output directory (default: reports)
- `--format`: Output format: json,markdown,html,all
- `--live-report`: Enable live report updates

#### Performance Tuning
- `--threads, -j`: Number of concurrent threads (default: 50)
- `--rate-limit, -l`: Requests per second limit (default: 150)
- `--timeout`: Overall scan timeout
- `--no-throttle`: Disable request throttling

#### Module-Specific Flags
- `--llm-model`: Path to LLM model file (.gguf)
- `--temp`: LLM temperature (default: 0.7)
- `--generate-patch`: Generate auto-patch recommendations
- `--shadow-clone`: Shadow clone proxy (format: port:9009)
- `--dep-drift`: Scan for dependency drift
- `--timing-map`: Perform side-channel timing analysis
- `--blue-team`: Generate blue team replay kit

## Output Formats

### JSON Output

Raw scan results in JSON format:

```bash
tracehawkx scan --target example.com --format json
```

Output: `results/scan-results.json`

### Markdown Report

Human-readable report in Markdown:

```bash
tracehawkx scan --target example.com --format markdown
```

Output: `reports/scan-report.md`

### HTML Report

Styled HTML report (requires pandoc):

```bash
tracehawkx scan --target example.com --format html
```

Output: `reports/scan-report.html`

### All Formats

Generate all output formats:

```bash
tracehawkx scan --target example.com --format all
```

## Module Management

### List All Modules

```bash
tracehawkx list-modules
```

### Module Categories

- **Stable**: Production-ready modules
- **Bleeding-Edge**: Experimental modules

### Module Prerequisites

Each module checks for required tools automatically. Install missing tools:

```bash
# Ubuntu/Debian
sudo apt install nmap nuclei subfinder httpx

# macOS
brew install nmap nuclei subfinder httpx

# Arch Linux
sudo pacman -S nmap
```

## Advanced Usage

### Multiple Targets

```bash
# Multiple domains
tracehawkx scan --target example.com,test.com

# CIDR range
tracehawkx scan --target 192.168.1.0/24

# From file
echo -e "example.com\ntest.com" > targets.txt
tracehawkx scan --target-file targets.txt
```

### Custom Configuration

```yaml
# ~/.tracehawkx.yaml
output: /tmp/scans
report_dir: /tmp/reports
threads: 100
rate_limit: 200
bleeding_edge: true
generate_patch: true
```

### Environment Variables

```bash
export TRACEHAWK_CONFIG=/path/to/config.yaml
export TRACEHAWK_OUTPUT=/path/to/output
```

### Docker Usage

```bash
# Basic scan
docker run --rm -v $(pwd)/results:/app/results tracehawkx scan --target example.com

# With custom config
docker run --rm \
  -v $(pwd)/config.yaml:/app/.tracehawkx.yaml \
  -v $(pwd)/results:/app/results \
  tracehawkx scan --target example.com
```

## Examples

### Basic Web Application Assessment

```bash
tracehawkx scan \
  --target webapp.example.com \
  --deep \
  --output webapp-scan \
  --report webapp-report
```

### Network Range Reconnaissance

```bash
tracehawkx scan \
  --target 10.0.0.0/24 \
  --exclude 10.0.0.1,10.0.0.254 \
  --stealth \
  --threads 20
```

### Advanced Security Assessment

```bash
tracehawkx scan \
  --target api.example.com \
  --bleeding-edge \
  --generate-patch \
  --llm-model /path/to/model.gguf \
  --blue-team
```

### Continuous Monitoring

```bash
#!/bin/bash
# monitor.sh
while true; do
  tracehawkx scan \
    --target-file targets.txt \
    --output "scans/$(date +%Y%m%d-%H%M%S)" \
    --format json
  sleep 3600  # Scan every hour
done
```

### Integration with CI/CD

```yaml
# .github/workflows/security-scan.yml
name: Security Scan
on:
  schedule:
    - cron: '0 2 * * *'  # Daily at 2 AM

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run TraceHawk X
        run: |
          docker run --rm \
            -v ${{ github.workspace }}/results:/app/results \
            tracehawkx scan --target ${{ secrets.TARGET_DOMAIN }}
      - name: Upload Results
        uses: actions/upload-artifact@v2
        with:
          name: security-scan-results
          path: results/
```

## Troubleshooting

### Common Issues

1. **Permission Denied**: Run with sudo for network namespace features
2. **Module Not Found**: Install required tools (nmap, nuclei, etc.)
3. **Rate Limiting**: Reduce `--rate-limit` or increase `--threads`
4. **Memory Usage**: Reduce `--threads` for large scans

### Debug Mode

```bash
tracehawkx scan --target example.com --debug
```

### Verbose Output

```bash
tracehawkx scan --target example.com --verbose
```

### Log Files

Logs are written to stderr by default. Redirect to file:

```bash
tracehawkx scan --target example.com 2> scan.log
``` 