# Simple Plugin Example Makefile

.PHONY: all build clean run test

# Variables
PLUGIN_NAME = plugin-hello
HOST_NAME = host
GO = go
GOFLAGS = -v

# Default target
all: build

# Build plugin and host
build: build-plugin build-host

# Build the plugin
build-plugin:
	@echo "Building plugin..."
	@mkdir -p ./plugins
	$(GO) build $(GOFLAGS) -o ./plugins/$(PLUGIN_NAME) ./plugin

# Build the host
build-host:
	@echo "Building host..."
	$(GO) build $(GOFLAGS) -o ./$(HOST_NAME) ./host

# Run the example
run: build
	@echo "Running example..."
	./$(HOST_NAME)

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf ./plugins
	@rm -f ./$(HOST_NAME)
	@rm -f *.log

# Development mode with hot reload
dev:
	@echo "Starting in development mode..."
	@# This would include file watching in a real implementation
	$(MAKE) run

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy