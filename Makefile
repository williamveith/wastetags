# Root Directories
BINARY_ROOT_DIR := bin
BUILD_DIR := build
BUILD_LOG_DIR := $(BINARY_ROOT_DIR)/logs

BINARY_NAME := wastetags
DOCKER_IMAGE := $(BINARY_NAME):latest
DOCKERFILE := $(BUILD_DIR)/Dockerfile
DATA_FILE := build/wastetags.sqlite3
USER := pi.local

# Proto files
PROTO_DIR = protobuf
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

all: clean dev run-dev

# Build for Linux
linux: BUILD_TYPE := linux
linus: TARGET_ARCH := arm64
linux: TARGET_VARIANT := v8
linux: BINARY_DIR := $(BINARY_ROOT_DIR)/$(BUILD_TYPE)
linux: BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)
linux:
	@echo "Building for Linux..."
	@mkdir -p $(BINARY_DIR)
	@docker build \
		--build-arg TARGETOS=$(BUILD_TYPE) \
		--build-arg TARGETARCH=$(TARGET_ARCH) \
		--build-arg TARGETVARIANT=$(TARGET_VARIANT) \
		-t $(DOCKER_IMAGE) \
		-f $(DOCKERFILE) .
	@docker run --rm -v $$(pwd)/$(BINARY_DIR):/export $(DOCKER_IMAGE) cp /$(BINARY_NAME) /export/
	@cp -r $(BUILD_DIR)/$(BUILD_TYPE)/* $(BINARY_DIR)/
	@echo "$(BUILD_TYPE) build complete. Files located in $$(pwd)/$(BINARY_DIR)"

push-linux: BUILD_TYPE := linux
push-linux: BINARY_DIR := $(BINARY_ROOT_DIR)/$(BUILD_TYPE)
push-linux: clean linux
	@scp -r $(BINARY_DIR) $(USER):/tmp/$(BUILD_TYPE)
	@scp $(DATA_FILE) $(USER):/tmp/$(BUILD_TYPE)
	@ssh $(USER) 'sudo mv /tmp/linux/wastetags /usr/local/bin/ && \
	sudo mv /tmp/linux/wastetags.sqlite3 /var/lib/wastetags/ && \
	sudo rm -rf /tmp/linux && \
	sudo systemctl restart wastetags'
	@echo "Push complete. Files moved to final destinations on remote server."

mac: BUILD_TYPE := darwin
mac: TARGET_ARCH := arm64
mac: TARGET_VARIANT := v8
mac: BINARY_DIR := $(BINARY_ROOT_DIR)/$(BUILD_TYPE)
mac: BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)
mac:
	@echo "Building for $(BUILD_TYPE)/$(TARGET_ARCH)/$(TARGET_VARIANT)..."
	@mkdir -p $(BINARY_DIR)
	@mkdir -p $(BUILD_LOG_DIR)
	@go build -v -gcflags="all=-N -l" -x -ldflags="-s -w" -o $(BINARY_PATH) cmd/$(BINARY_NAME)/*.go > $(BUILD_LOG_DIR)/$(BUILD_TYPE).log 2>&1
	@cp $(BUILD_DIR)/$(BUILD_TYPE)/iconfile.icns $(BINARY_DIR)/iconfile.icns
	@cp $(BUILD_DIR)/$(BUILD_TYPE)/Info.plist $(BINARY_DIR)/Info.plist
	@echo "$(BUILD_TYPE) build complete. Files located in $(BINARY_DIR)"

# Build for dev
dev: BUILD_TYPE := dev
dev: BINARY_DIR := $(BINARY_ROOT_DIR)/$(BUILD_TYPE)
dev: BINARY_PATH := $(BINARY_DIR)/$(BINARY_NAME)
dev:
	@echo "Building for dev..."
	@mkdir -p $(BINARY_DIR)
	@mkdir -p $(BUILD_LOG_DIR)
	@go build -v -gcflags="all=-N -l" -x -ldflags="-s -w" -o $(BINARY_PATH) cmd/$(BINARY_NAME)/*.go > $(BUILD_LOG_DIR)/$(BUILD_TYPE).log 2>&1
	@echo "$(BUILD_TYPE) build complete. Files located in $(BINARY_DIR)"

# Run the application for a specific build
run-linux: BUILD_TYPE := linux
run-linux:
	@echo "Running $(BUILD_TYPE) build..."
	@$(BINARY_ROOT_DIR)/$(BUILD_TYPE)/$(BINARY_NAME) --config=$(BINARY_ROOT_DIR)/$(BUILD_TYPE)/config.json

run-mac: BUILD_TYPE := darwin
run-mac:
	@echo "Running $(BUILD_TYPE) build..."
	@$(BINARY_ROOT_DIR)/$(BUILD_TYPE)/$(BINARY_NAME) --config=$(BINARY_ROOT_DIR)/$(BUILD_TYPE)/config.json

run-dev: BUILD_TYPE := dev
run-dev:
	@echo "Running $(BUILD_TYPE) build at http://localhost:8080"
	@$(BINARY_ROOT_DIR)/$(BUILD_TYPE)/$(BINARY_NAME) --$(BUILD_TYPE)

proto:
	@protoc --go_out=. $(PROTO_FILES)

# Clean build artifacts
clean:
	@echo "Cleaning all build artifacts..."
	@rm -rf $(BINARY_ROOT_DIR)
	@rm -rf $(BUILD_LOG_DIR)
	@echo "Clean complete."

mega-git:
	@cp -rf cmd/wastetags/data ~/Projects/MegaGithub/wastetags/
	@git --git-dir=$$HOME/Projects/MegaGithub/wastetags/.git --work-tree=$$HOME/Projects/MegaGithub/wastetags add .
	@git --git-dir=$$HOME/Projects/MegaGithub/wastetags/.git --work-tree=$$HOME/Projects/MegaGithub/wastetags commit -m "Data update"
	@mega-put /Users/main/Projects/MegaGithub/wastetags/.git /github/wastetags/.git


# Help command to display available targets
help:
	@echo "\n"
	@echo "Makefile commands:"
	@echo "  make linux         - Builds the project for Linux"
	@echo "  make mac           - Builds the project for macOS"
	@echo "  make dev           - Builds the project for dev"
	@echo "  make run-linux     - Runs the Linux build"
	@echo "  make run-mac       - Runs the macOS build"
	@echo "  make run-dev       - Runs the dev build"
	@echo "  make clean         - Cleans all build artifacts"
	@echo "  make help          - Displays this help message"
	@echo "\n"

.PHONY: all linux mac dev run-linux run-mac run-dev clean help
