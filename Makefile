# Root Directories
BINARY_ROOT_DIR := bin
CONFIG_DIR := configs
BUILD_DIR := build

BINARY_NAME := wastetags
DOCKER_IMAGE := $(BINARY_NAME):latest
DOCKERFILE := $(BUILD_DIR)/Dockerfile
DATA_FILE := data/chemicals.sqlite3
USER := pi.local

# Default target
all: linux

# Build for Linux
linux: BUILD_TYPE := linux
linux: BINARY_DIR := $(BINARY_ROOT_DIR)/$(BUILD_TYPE)
linux: BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)
linux:
	@echo "Building for Linux..."
	@mkdir -p $(BINARY_DIR)
	@docker build \
		--build-arg TARGETPLATFORM=$(BUILD_TYPE)/arm64 \
		--build-arg TARGETOS=$(BUILD_TYPE) \
		--build-arg TARGETARCH=arm64 \
		--build-arg TARGETVARIANT=v8 \
		-t $(DOCKER_IMAGE) \
		-f $(DOCKERFILE) .
	@docker run --rm -v $$(pwd)/$(BINARY_DIR):/export $(DOCKER_IMAGE) cp /wastetags /export/
	@cp $(CONFIG_DIR)/$(BUILD_TYPE).json $(BINARY_DIR)/config.json
	@cp $(BUILD_DIR)/$(BUILD_TYPE)/wastetags.service $(BINARY_DIR)/wastetags.service
	@echo "$(BUILD_TYPE) build complete. Files located in $(BINARY_DIR)"

push-linux: clean linux
	@scp -r $(BINARY_ROOT_DIR)/linux $(USER):/tmp/linux
	@scp $(DATA_FILE) $(USER):/tmp/linux
	@ssh $(USER) 'sudo mv /tmp/linux/wastetags /usr/local/bin/ && \
	sudo mv /tmp/linux/config.json /etc/wastetags/ && \
	sudo mv /tmp/linux/chemicals.sqlite3 /var/lib/wastetags/ && \
	sudo rm -rf /tmp/linux && \
	sudo systemctl restart wastetags'
	@echo "Push complete. Files moved to final destinations on remote server."

# Build for macOS
mac: BUILD_TYPE := macos
mac: BINARY_PATH := $(BINARY_ROOT_DIR)/macos/$(BINARY_NAME)
mac:
	@echo "Building for macOS..."
	@mkdir -p $(BINARY_ROOT_DIR)/macos
	@go build -ldflags="-s -w" -o $(BINARY_PATH) cmd/$(BINARY_NAME)/*.go
	@cp $(CONFIG_DIR)/macos.json $(BINARY_ROOT_DIR)/macos/config.json
	@cp $(DATA_FILE) $(BINARY_ROOT_DIR)/macos/chemicals.sqlite3
	@echo "macOS build complete. Files located in $(BINARY_ROOT_DIR)/macos"

# Build for dev
dev: BUILD_TYPE := dev
dev: BINARY_DIR := $(BINARY_ROOT_DIR)/$(BUILD_TYPE)
dev: BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)
dev:
	@echo "Building for dev..."
	@mkdir -p $(BINARY_DIR)
	@go build -ldflags="-s -w" -o $(BINARY_PATH) cmd/$(BINARY_NAME)/*.go
	@cp $(CONFIG_DIR)/$(BUILD_TYPE).json $(BINARY_DIR)/config.json
	@echo "$(BUILD_TYPE) build complete. Files located in $(BINARY_DIR)"

# Run the application for a specific build
run-linux:
	@echo "Running Linux build..."
	@$(BINARY_ROOT_DIR)/linux/$(BINARY_NAME) --config=$(BINARY_ROOT_DIR)/linux/config.json

run-mac:
	@echo "Running macOS build..."
	@$(BINARY_ROOT_DIR)/macos/$(BINARY_NAME) --config=$(BINARY_ROOT_DIR)/macos/config.json

run-dev:
	@echo "Running dev build..."
	@$(BINARY_ROOT_DIR)/dev/$(BINARY_NAME) --config=$(BINARY_ROOT_DIR)/dev/config.json

# Clean build artifacts
clean:
	@echo "Cleaning all build artifacts..."
	@rm -rf $(BINARY_ROOT_DIR)
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
