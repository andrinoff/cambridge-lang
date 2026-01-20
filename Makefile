# Cambridge Pseudocode Interpreter Makefile
# Based on Cambridge International AS & A Level Computer Science 9618

BINARY_NAME=cambridge
VERSION=1.0.0
BUILD_DIR=build
GO=go
GOFLAGS=-ldflags="-s -w"

# Platform detection
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) ./cmd/cambridge

# Build the Language Server Protocol
.PHONY: build-lsp
build-lsp:
	$(GO) build $(GOFLAGS) -o cambridge-lsp ./cmd/cambridge-lsp

# Build with debug symbols
.PHONY: build-debug
build-debug:
	$(GO) build -o $(BINARY_NAME) ./cmd/cambridge

# Run tests
.PHONY: test
test:
	$(GO) test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -f examples/output.txt

# Install to GOPATH/bin
.PHONY: install
install:
	$(GO) install ./cmd/cambridge

# Uninstall from GOPATH/bin
.PHONY: uninstall
uninstall:
	rm -f $(shell go env GOPATH)/bin/$(BINARY_NAME)

# Format code
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# Lint code
.PHONY: lint
lint:
	$(GO) vet ./...

# Run all examples
.PHONY: examples
examples: build
	@echo "Running all examples..."
	@echo ""
	@echo "=== hello.csal ==="
	./$(BINARY_NAME) run examples/hello.csal
	@echo ""
	@echo "=== variables.csal ==="
	./$(BINARY_NAME) run examples/variables.csal
	@echo ""
	@echo "=== selection.csal ==="
	./$(BINARY_NAME) run examples/selection.csal
	@echo ""
	@echo "=== loops.csal ==="
	./$(BINARY_NAME) run examples/loops.csal
	@echo ""
	@echo "=== functions.csal ==="
	./$(BINARY_NAME) run examples/functions.csal
	@echo ""
	@echo "=== arrays.csal ==="
	./$(BINARY_NAME) run examples/arrays.csal
	@echo ""
	@echo "=== strings.csal ==="
	./$(BINARY_NAME) run examples/strings.csal
	@echo ""
	@echo "=== records.csal ==="
	./$(BINARY_NAME) run examples/records.csal

# Start REPL
.PHONY: repl
repl: build
	./$(BINARY_NAME) repl

# Cross-compilation targets
.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: build-linux
build-linux:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/cambridge
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/cambridge

.PHONY: build-darwin
build-darwin:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/cambridge
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/cambridge

.PHONY: build-windows
build-windows:
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/cambridge

# Create release archives
.PHONY: release
release: build-all
	@mkdir -p $(BUILD_DIR)/release
	cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	cd $(BUILD_DIR) && zip release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe

# Show help
.PHONY: help
help:
	@echo "Cambridge Pseudocode Interpreter - Makefile targets"
	@echo ""
	@echo "Build targets:"
	@echo "  make build        - Build the binary"
	@echo "  make build-debug  - Build with debug symbols"
	@echo "  make build-all    - Cross-compile for all platforms"
	@echo "  make install      - Install to GOPATH/bin"
	@echo "  make uninstall    - Remove from GOPATH/bin"
	@echo ""
	@echo "Development targets:"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage- Run tests with coverage report"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Lint code"
	@echo "  make clean        - Clean build artifacts"
	@echo ""
	@echo "Run targets:"
	@echo "  make repl         - Start interactive REPL"
	@echo "  make examples     - Run all example programs"
	@echo ""
	@echo "Release targets:"
	@echo "  make build-linux  - Build for Linux (amd64, arm64)"
	@echo "  make build-darwin - Build for macOS (amd64, arm64)"
	@echo "  make build-windows- Build for Windows (amd64)"
	@echo "  make release      - Create release archives"
