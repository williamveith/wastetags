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

# scp /Users/main/Projects/Go/wastetags/bin/pi/wastetags pi.local:~
# scp /Users/main/Projects/Go/wastetags/bin/pi/config.json pi.local:~
# scp /Users/main/Projects/Go/wastetags/bin/pi/chemicals.sqlite3 pi.local:~
# scp /Users/main/Projects/Go/wastetags/bin/pi/wastetags.service pi.local:~
pi: clean
	@echo "Building for Raspberry Pi..."
	@mkdir -p $$(pwd)/bin/pi
	@docker build -t wastetags:latest -f build/pi/Dockerfile .
	@docker run --rm -v $$(pwd)/bin/pi:/export wastetags:latest cp /wastetags /export/
	@cp $$(pwd)/data/chemicals.sqlite3 $$(pwd)/bin/pi/chemicals.sqlite3
	@cp $$(pwd)/build/pi/wastetags.service $$(pwd)/bin/pi/wastetags.service
	@cp $$(pwd)/build/pi/config.json $$(pwd)/bin/pi/config.json

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
