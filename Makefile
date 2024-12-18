# Variables
BINARY_NAME := wastetags
BINARY_DIR := bin
BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)
CMD_DIR := cmd/$(BINARY_NAME)
SRC_FILES := $(wildcard $(CMD_DIR)/*.go)
DOCKER_COMPOSE_DIR := deployments
DOCKER_COMPOSE_FILE := $(DOCKER_COMPOSE_DIR)/docker-compose.yml

# Default target
all: build

# Build the binary
build: clean
	@echo "Building the project..."
	@mkdir -p $(BINARY_DIR)
	@go build -ldflags="-s -w" -o $(BINARY_PATH) $(SRC_FILES)
	@echo "Build complete. Binary located at $(BINARY_PATH)"

pi: clean
	@echo "Building for Raspberry Pi..."
	@mkdir -p $$(pwd)/bin/pi/data
	@docker build -t wastetags:latest -f build/package/Dockerfile .
	@docker run --rm -v $$(pwd)/bin/pi:/export wastetags:latest cp /wastetags /export/
	@cp $$(pwd)/data/chemicals.sqlite3 $$(pwd)/bin/pi/data/chemicals.sqlite3

# Run the application
run: build
	@echo "Running the application..."
	@$(BINARY_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINARY_DIR)
	@mkdir -p $(BINARY_DIR)
	@echo "Clean complete."

# Help command to display available targets
help:
	@echo "Makefile commands:"
	@echo "  make              - Builds the project"
	@echo "  make build        - Builds the project"
	@echo "  make run          - Runs the application"
	@echo "  make clean        - Cleans build artifacts"
	@echo "  make help         - Displays this help message"

.PHONY: all build run clean test deps fmt lint generate docker-up docker-down docker-build-no-cache docker-logs help
