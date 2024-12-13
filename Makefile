# Variables
BINARY_NAME := wastetags
BINARY_DIR := bin
BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)
CMD_DIR := cmd/$(BINARY_NAME)
DOCKER_COMPOSE_DIR := deployments
DOCKER_COMPOSE_FILE := $(DOCKER_COMPOSE_DIR)/docker-compose.yml

# Default target
all: build

# Build the binary
build: clean
	@echo "Building the project..."
	@mkdir -p $(BINARY_DIR)
	@go build -ldflags="-s -w" -o $(BINARY_PATH) $(CMD_DIR)/main.go
	@echo "Build complete. Binary located at $(BINARY_PATH)"

# Run the application
run: build
	@echo "Running the application..."
	@$(BINARY_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINARY_DIR)
	@echo "Clean complete."

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Install dependencies
deps:
	@echo "Ensuring dependencies are up to date..."
	@go mod tidy

# Format the code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint the code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Generate code (if applicable)
generate:
	@echo "Generating code..."
	@go generate ./...

# Build and run Docker Compose services
docker-up:
	@echo "Starting Docker Compose services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d

# Stop Docker Compose services
docker-down:
	@echo "Stopping Docker Compose services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Build Docker image without cache
docker-build-no-cache:
	@echo "Building Docker image without cache..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build --no-cache

# View logs
docker-logs:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Help command to display available targets
help:
	@echo "Makefile commands:"
	@echo "  make              - Builds the project"
	@echo "  make build        - Builds the project"
	@echo "  make run          - Runs the application"
	@echo "  make clean        - Cleans build artifacts"
	@echo "  make test         - Runs tests"
	@echo "  make deps         - Installs dependencies"
	@echo "  make fmt          - Formats the code"
	@echo "  make lint         - Lints the code (requires golangci-lint)"
	@echo "  make generate     - Runs code generation"
	@echo "  make docker-up    - Builds and starts Docker Compose services"
	@echo "  make docker-down  - Stops Docker Compose services"
	@echo "  make docker-build-no-cache - Builds Docker image without cache"
	@echo "  make docker-logs  - Follows logs from Docker Compose services"
	@echo "  make help         - Displays this help message"

.PHONY: all build run clean test deps fmt lint generate docker-up docker-down docker-build-no-cache docker-logs help
