# TraceHawk X Makefile

# Build variables
BINARY_NAME=tracehawkx
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Directories
BUILD_DIR=build
DIST_DIR=dist
CMD_DIR=cmd/tracehawkx

.PHONY: all build clean test deps lint install uninstall release docker help

# Default target
all: clean deps build

# Build the binary
build:
	@echo "🔨 Building TraceHawk X..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies updated"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...
	@echo "✅ Tests complete"

# Run linter
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found, skipping lint"; \
	fi

# Install binary to system
install: build
	@echo "📥 Installing TraceHawk X..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"

# Uninstall binary from system
uninstall:
	@echo "🗑️  Uninstalling TraceHawk X..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Uninstalled"

# Build release binaries for multiple platforms
release: clean deps
	@echo "🚀 Building release binaries..."
	@mkdir -p $(DIST_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	
	@echo "✅ Release binaries built in $(DIST_DIR)/"

# Build Docker image
docker:
	@echo "🐳 Building Docker image..."
	docker build -t tracehawkx:$(VERSION) .
	docker tag tracehawkx:$(VERSION) tracehawkx:latest
	@echo "✅ Docker image built: tracehawkx:$(VERSION)"

# Development setup
dev-setup:
	@echo "🛠️  Setting up development environment..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@echo "✅ Development environment ready"

# Run the binary with example target
run: build
	@echo "🏃 Running TraceHawk X..."
	./$(BUILD_DIR)/$(BINARY_NAME) --help

# Show help
help:
	@echo "TraceHawk X - Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  all        - Clean, download deps, and build"
	@echo "  build      - Build the binary"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  install    - Install binary to /usr/local/bin"
	@echo "  uninstall  - Remove binary from /usr/local/bin"
	@echo "  release    - Build release binaries for all platforms"
	@echo "  docker     - Build Docker image"
	@echo "  dev-setup  - Set up development environment"
	@echo "  run        - Build and run with --help"
	@echo "  help       - Show this help message"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION    - Version string (default: dev)"
	@echo "  COMMIT     - Git commit hash (auto-detected)"
	@echo "  DATE       - Build date (auto-generated)" 