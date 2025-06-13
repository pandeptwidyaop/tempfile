# TempFiles - Makefile

.PHONY: build run dev test clean help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_NAME=tempfile
BINARY_UNIX=$(BINARY_NAME)_unix

# Build directory
BUILD_DIR=./bin

# Main package path
MAIN_PATH=./cmd/server

# Default target
all: test build

## Build the application
build:
	@echo "🔨 Building application..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## Run the application in development mode
dev:
	@echo "🚀 Starting development server..."
	$(GOCMD) run $(MAIN_PATH)

## Run the built application
run: build
	@echo "🚀 Starting application..."
	./$(BUILD_DIR)/$(BINARY_NAME)

## Test the application
test:
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

## Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

## Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies updated"

## Build for Linux
build-linux:
	@echo "🐧 Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v $(MAIN_PATH)
	@echo "✅ Linux build complete: $(BUILD_DIR)/$(BINARY_UNIX)"

## Docker build
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t tempfile:latest .
	@echo "✅ Docker build complete"

## Docker run
docker-run:
	@echo "🐳 Running Docker container..."
	docker run -p 3000:3000 tempfile:latest

## Docker compose up
docker-up:
	@echo "🐳 Starting with Docker Compose..."
	docker-compose up -d

## Docker compose down
docker-down:
	@echo "🐳 Stopping Docker Compose..."
	docker-compose down

## Test Docker health check
docker-health:
	@echo "🏥 Testing Docker health check..."
	curl -f http://localhost:3000/health || echo "Health check failed"

## Show help
help:
	@echo "📋 Available commands:"
	@echo ""
	@echo "  build        - Build the application"
	@echo "  dev          - Run in development mode"
	@echo "  run          - Build and run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  build-linux  - Build for Linux"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-up    - Start with Docker Compose"
	@echo "  docker-down  - Stop Docker Compose"
	@echo "  docker-health - Test Docker health check"
	@echo "  help         - Show this help"
	@echo ""
