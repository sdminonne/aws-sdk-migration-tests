# Makefile for aws-sdk-migration-tests

# Binary names
CROSS_VERSION_BIN := cross_version_infrastructure
MIXED_SDK_BIN := mixed_sdk

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get

# Build flags
LDFLAGS := -ldflags="-s -w"

.PHONY: all cross_version mixed_sdk clean test

# Default target - build all binaries
all: cross_version mixed_sdk

# Build cross_version_infrastructure binary
cross_version:
	$(GOBUILD) $(LDFLAGS) -o $(CROSS_VERSION_BIN) cross_version_infrastructure.go

# Build mixed_sdk binary
mixed_sdk:
	$(GOBUILD) $(LDFLAGS) -o $(MIXED_SDK_BIN) mixed_sdk.go

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(CROSS_VERSION_BIN)
	rm -f $(MIXED_SDK_BIN)

# Display help information
help:
	@echo "Available targets:"
	@echo "  all            - Build both binaries (default)"
	@echo "  cross_version  - Build cross_version_infrastructure binary"
	@echo "  mixed_sdk      - Build mixed_sdk binary"
	@echo "  test           - Run tests"
	@echo "  clean          - Remove built binaries"
	@echo "  help           - Display this help message"
