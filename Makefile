# Variables
BINARY_NAME := wastetags
BINARY_DIR := bin
CONFIG_DIR := configs
DATA_FILE := data/chemicals.sqlite3

# Default target
all: linux

# Build for Linux
linux: BUILD_TYPE := linux
linux: BINARY_PATH := $(BINARY_DIR)/linux/$(BINARY_NAME)
linux:
	@echo "Building for Linux..."
	@mkdir -p $(BINARY_DIR)/linux
	@docker build -t wastetags:latest -f build/linux/Dockerfile .
	@docker run --rm -v $$(pwd)/bin/linux:/export wastetags:latest cp /wastetags /export/
	@cp $(CONFIG_DIR)/linux.json $(BINARY_DIR)/linux/config.json
	@cp $(DATA_FILE) $(BINARY_DIR)/linux/chemicals.sqlite3
	@echo "Linux build complete. Files located in $(BINARY_DIR)/linux"

# Build for macOS
mac: BUILD_TYPE := macos
mac: BINARY_PATH := $(BINARY_DIR)/macos/$(BINARY_NAME)
mac:
	@echo "Building for macOS..."
	@mkdir -p $(BINARY_DIR)/macos
	@go build -ldflags="-s -w" -o $(BINARY_PATH) cmd/$(BINARY_NAME)/*.go
	@cp $(CONFIG_DIR)/macos.json $(BINARY_DIR)/macos/config.json
	@cp $(DATA_FILE) $(BINARY_DIR)/macos/chemicals.sqlite3
	@echo "macOS build complete. Files located in $(BINARY_DIR)/macos"

# Build for dev
dev: BUILD_TYPE := dev
dev: BINARY_PATH := $(BINARY_DIR)/dev/$(BINARY_NAME)
dev:
	@echo "Building for dev..."
	@mkdir -p $(BINARY_DIR)/dev
	@go build -ldflags="-s -w" -o $(BINARY_PATH) cmd/$(BINARY_NAME)/*.go
	@cp $(CONFIG_DIR)/dev.json $(BINARY_DIR)/dev/config.json
	@echo "dev build complete. Files located in $(BINARY_DIR)/dev"

# Run the application for a specific build
run-linux:
	@echo "Running Linux build..."
	@$(BINARY_DIR)/linux/$(BINARY_NAME) --config=$(BINARY_DIR)/linux/config.json

run-mac:
	@echo "Running macOS build..."
	@$(BINARY_DIR)/macos/$(BINARY_NAME) --config=$(BINARY_DIR)/macos/config.json

run-dev:
	@echo "Running dev build..."
	@$(BINARY_DIR)/dev/$(BINARY_NAME) --config=$(BINARY_DIR)/dev/config.json

# Clean build artifacts
clean:
	@echo "Cleaning all build artifacts..."
	@rm -rf $(BINARY_DIR)
	@echo "Clean complete."

# Help command to display available targets
help:
	@echo "Makefile commands:"
	@echo "  make linux         - Builds the project for Linux"
	@echo "  make mac           - Builds the project for macOS"
	@echo "  make dev       	- Builds the project for dev"
	@echo "  make run-linux     - Runs the Linux build"
	@echo "  make run-mac       - Runs the macOS build"
	@echo "  make run-dev  		- Runs the dev build"
	@echo "  make clean         - Cleans all build artifacts"
	@echo "  make help          - Displays this help message"

.PHONY: all linux mac dev run-linux run-mac run-dev clean help
