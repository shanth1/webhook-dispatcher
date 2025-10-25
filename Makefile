# Makefile for the github-webhook-dispatcher project

# --- Variables ---

# Main binary name (complete version)
BINARY_NAME_COMPLETE := gitrelay
# Simple binary name (basic version)
BINARY_NAME_BASIC := basic-webhook-handler

# Paths to source code
CMD_PATH_COMPLETE := ./cmd/complete
CMD_PATH_BASIC := ./cmd/basic

# Config file for running the application
CONFIG_FILE := config/example.yaml

# List of all Go files to track for changes
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*")


# --- Main Targets ---

# Default target: build the main application
.PHONY: all
all: build

# Display help for the available commands
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make all                Build the main application (default)."
	@echo "  make build              Build the main application (gitrelay)."
	@echo "  make build-basic        Build the basic application (basic-webhook-handler)."
	@echo "  make build-all          Build both applications."
	@echo ""
	@echo "  make run                Build and run the main application with config/example.yaml."
	@echo "  make run-basic          Build and run the basic application."
	@echo ""
	@echo "  make test               Run Go tests for the entire project."
	@echo "  make test-webhook       Send all test webhooks using webhook.sh."
	@echo "  make test-webhook       EVENT=<event_name> Send a specific webhook (e.g., EVENT=push)."
	@echo "  make test-push          Shortcut for 'make test-webhook EVENT=push'."
	@echo "  make test-ping          Shortcut for 'make test-webhook EVENT=ping'."
	@echo ""
	@echo "  make clean              Remove built binaries."
	@echo "  make tidy               Run 'go mod tidy'."
	@echo "  make fmt                Format Go code using 'go fmt'."
	@echo "  make vet                Run 'go vet' to check for issues."
	@echo "  make lint               Run the golangci-lint linter (if installed)."


# --- Build ---

# Build the complete application version
.PHONY: build-complete
build-complete: $(GOFILES)
	@echo "Building $(BINARY_NAME_COMPLETE)..."
	@go build -o build/$(BINARY_NAME_COMPLETE) $(CMD_PATH_COMPLETE)
	@echo "Build of build/$(BINARY_NAME_COMPLETE) complete."

build-complete-linux: $(GOFILES)
	@echo "Building $(BINARY_NAME_COMPLETE) for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o build/$(BINARY_NAME_COMPLETE) $(CMD_PATH_COMPLETE)
	@echo "Build of build/$(BINARY_NAME_COMPLETE) complete."

# Build the basic application version
.PHONY: build-basic
build-basic: $(GOFILES)
	@echo "Building $(BINARY_NAME_BASIC)..."
	@go build -o build/$(BINARY_NAME_BASIC) $(CMD_PATH_BASIC)
	@echo "Build of build/$(BINARY_NAME_BASIC) complete."

# Build both applications
.PHONY: build-all
build-all: build-complete build-basic


# --- Run ---

# Run the main application
.PHONY: run
run: build-complete
	@echo "Running $(BINARY_NAME_COMPLETE) with config $(CONFIG_FILE)..."
	@./build/$(BINARY_NAME_COMPLETE) --config $(CONFIG_FILE)

# Run the basic application
.PHONY: run-basic
run-basic: build-basic
	@echo "Running $(BINARY_NAME_BASIC)..."
	@./build/$(BINARY_NAME_BASIC)


# --- Testing ---

# Run standard Go tests
.PHONY: test
test:
	@echo "Running Go tests..."
	@go test -v ./...

# Send test webhooks using webhook.sh
# Specify an event: make test-webhook EVENT=push
.PHONY: test-webhook
test-webhook:
	# Ensure the script is executable
	@chmod +x ./webhook.sh
	# Set a default value for EVENT if not provided
	$(eval EVENT ?= all)
	@echo "Sending test webhook event: $(EVENT)..."
	@if [ "$(EVENT)" = "all" ]; then \
		./webhook.sh; \
	else \
		./webhook.sh $(EVENT); \
	fi

# Convenience shortcuts for sending specific events
.PHONY: test-push
test-push: ; $(MAKE) test-webhook EVENT=push

.PHONY: test-ping
test-ping: ; $(MAKE) test-webhook EVENT=ping

.PHONY: test-issues
test-issues: ; $(MAKE) test-webhook EVENT=issues


# --- Maintenance ---

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up binaries..."
	@rm -rf build/

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Running go mod tidy..."
	@go mod tidy

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code for issues
.PHONY: vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Lint with golangci-lint
.PHONY: lint
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo >&2 "golangci-lint is not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	@echo "Running golangci-lint linter..."
	@golangci-lint run ./...
