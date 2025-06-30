.PHONY: build build-arm32 build-arm64 build-all clean

# Common variables
pkg = github.com/padiazg/environment-monitor-daemon/models/version
ldflags = -X $(pkg).version=$(shell git describe --tags --always --dirty) 
ldflags += -X $(pkg).commit=$(shell git rev-parse HEAD)
ldflags += -X $(pkg).buildDate=$(shell date -Iseconds)
ldflags += -X $(pkg).buildBy=make

# Binary name (adjust as needed)
BINARY_NAME = environment-monitor-daemon

# Default build (current architecture)
build:
	@echo "Building $(BINARY_NAME) for current architecture..."
	@go build -ldflags "$(ldflags)" -o $(BINARY_NAME)

# ARM32 build (ARMv6 compatible - works on Raspberry Pi Zero/1)
build-arm32:
	@echo "Building $(BINARY_NAME) for ARM32 (ARMv6)..."
	@GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "$(ldflags)" -o $(BINARY_NAME)-arm32

# ARM64 build
build-arm64:
	@echo "Building $(BINARY_NAME) for ARM64..."
	@GOOS=linux GOARCH=arm64 go build -ldflags "$(ldflags)" -o $(BINARY_NAME)-arm64

# Build all architectures
build-all: build build-arm32 build-arm64
	@echo "All builds completed!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-arm32 $(BINARY_NAME)-arm64