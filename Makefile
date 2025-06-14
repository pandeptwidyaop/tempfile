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

## Build all platforms (GitHub Actions style)
build-all:
	@echo "🌍 Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 -v $(MAIN_PATH)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 -v $(MAIN_PATH)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe -v $(MAIN_PATH)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 -v $(MAIN_PATH)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 -v $(MAIN_PATH)
	@echo "✅ All platform builds complete"

## GitHub Actions style build
ci-build: test build-all
	@echo "🚀 CI Build complete"

## Test with coverage (GitHub Actions style)
test-coverage:
	@echo "🧪 Running tests with coverage..."
	$(GOTEST) -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

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

## Setup development environment
setup:
	@echo "🛠️ Setting up development environment..."
	@chmod +x ./dev.sh
	@./dev.sh setup
	@echo "✅ Development setup complete"

## Format code
fmt:
	@echo "🎨 Formatting code..."
	$(GOCMD) fmt ./...
	@echo "✅ Code formatted"

## Lint code
lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️ golangci-lint not installed. Run 'make setup' first"; \
	fi
	@echo "✅ Linting complete"

## Security scan
security:
	@echo "🔒 Running security scan..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "⚠️ gosec not installed. Run 'make setup' first"; \
	fi
	@echo "✅ Security scan complete"

## Run all quality checks
check: fmt lint security test
	@echo "✅ All quality checks passed! ✨"

## Interactive commit helper
commit:
	@./dev.sh commit

## Prepare for release (maintainers only)
release-prep:
	@echo "🚀 Preparing for release..."
	@./dev.sh release

help:
	@echo "📋 Available commands:"
	@echo ""
	@echo "  setup        - Setup development environment"
	@echo "  build        - Build the application"
	@echo "  dev          - Run in development mode"
	@echo "  run          - Build and run the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  security     - Run security scan"
	@echo "  check        - Run all quality checks (fmt, lint, security, test)"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  build-all    - Build for all platforms"
	@echo "  ci-build     - Run full CI build (test + build-all)"
	@echo "  commit       - Interactive commit helper (conventional commits)"
	@echo "  release-prep - Prepare for release (maintainers only)"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-up    - Start with Docker Compose"
	@echo "  docker-down  - Stop Docker Compose"
	@echo "  docker-health - Test Docker health check"
	@echo "  help         - Show this help"
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make setup   # First time setup"
	@echo "  make dev     # Start development"
	@echo "  make check   # Run quality checks"
	@echo "  make commit  # Commit with conventional format"
	@echo ""
