# Makefile for TurboGin Scaffolding

# Project variables
PROJECT_NAME := TurboGin
BINARY_NAME := turbogin
VERSION := 0.0.1
DOCKER_IMAGE := turbogin-app
DOCKER_TAG := latest

# Go variables
GO := go
GO_MODULE := TurboGin
WIRE := wire

# Directories
CONFIG_DIR := config
MIGRATIONS_DIR := migrations

# Default target
all: build

## -- Setup & Dependencies --
.PHONY: init
init:
	@echo "Initializing project..."
	$(GO) mod tidy
	$(GO) install github.com/google/wire/cmd/wire@latest

.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download

.PHONY: install-tools
install-tools:
	@echo "Installing required tools..."
	$(GO) install github.com/google/wire/cmd/wire@latest

.PHONY: generate
generate:
	@echo "Generating Wire dependencies..."
	$(WIRE) gen ./internal/wire

## -- Build & Run --
.PHONY: build
build: generate
	@echo "Building application with optimization..."
	$(GO) build -ldflags="-s -w" -trimpath -o bin/$(BINARY_NAME) ./cmd/server/main.go


.PHONY: build-linux
build-linux: generate
	@echo "Building Linux (amd64) application with optimization..."
	GOOS=linux GOARCH=amd64 $(GO) build \
		-ldflags="-s -w" \
		-trimpath \
		-o bin/$(BINARY_NAME)-linux ./cmd/server/main.go

.PHONY: run
run: generate
	@echo "Starting application..."
	$(GO) run cmd/server/main.go

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -rf docs/

## -- Docker --
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --name $(BINARY_NAME) --rm $(DOCKER_IMAGE):$(DOCKER_TAG)

## -- Help --
.PHONY: help
help:
	@echo "$(PROJECT_NAME) Scaffolding - Makefile Help"
	@echo ""
	@echo "Targets:"
	@echo "  init           - Initialize the project (go mod init)"
	@echo "  deps           - Download all dependencies"
	@echo "  install-tools  - Install required tools (wire)"
	@echo "  generate       - Generate Wire dependencies"
	@echo "  build          - Build the application"
	@echo "  build-linux    - Build the application for linux"
	@echo "  run            - Run the application"
	@echo "  clean          - Clean generated files"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  help           - Show this help message"