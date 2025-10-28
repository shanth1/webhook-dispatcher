# --- Variables ---

CMD_PATH := ./cmd
BINARY_NAME := hookrelay

CONFIG_FILE := config/example.yaml

# List of all Go files to track for changes
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*")


# --- Main Targets ---

# Display help for the available commands
.PHONY: help
help:
	@echo "Available commands:"
	@echo ""
	@echo "  make run                Build and run the main application with config/example.yaml."
	@echo "  make build              Build the main application (gitrelay)."
	@echo ""
	@echo "  make test               Run Go tests for the entire project."
	@echo "  make clean              Remove built binaries."


# --- Run ---

# Run the main application
.PHONY: run
run:
	@go run cmd/main.go --config $(CONFIG_FILE)


# --- Build ---

# Build the complete application version
.PHONY: build
build: $(GOFILES)
	@echo "Building $(BINARY_NAME)..."
	@go build -o build/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build of build/$(BINARY_NAME) complete."

build-linux: $(GOFILES)
	@echo "Building $(BINARY_NAME) for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o build/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build of build/$(BINARY_NAME) complete."


# --- Testing ---

# Run standard Go tests
.PHONY: test
test:
	@echo "Running Go tests..."
	@go test -v ./...


# --- Maintenance ---

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up binaries..."
	@rm -rf build/
