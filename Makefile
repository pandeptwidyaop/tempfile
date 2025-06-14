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

## Test rate limiting specifically
test-ratelimit:
	@echo "🧪 Running rate limiting tests..."
	$(GOTEST) -v ./internal/ratelimit/... ./internal/middleware/...

## Test with Redis (requires Redis running)
test-redis:
	@echo "🧪 Running Redis integration tests..."
	@if ! docker ps | grep -q redis-test; then \
		echo "🐳 Starting Redis test container..."; \
		docker run -d --name redis-test -p 6380:6379 redis:7-alpine; \
		sleep 2; \
	fi
	REDIS_URL=redis://localhost:6380 $(GOTEST) -v ./internal/ratelimit/redis_test.go
	@echo "🧹 Cleaning up Redis test container..."
	@docker stop redis-test && docker rm redis-test

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

## Docker GHCR Management (requires ./docker-manage.sh)
.PHONY: docker-pull docker-run-ghcr docker-stop docker-logs docker-health docker-clean docker-update docker-dev-ghcr docker-prod-ghcr

docker-pull:
	@echo "🐳 Pulling latest Docker image from GHCR..."
	@if [ -f ./docker-manage.sh ]; then \
		chmod +x ./docker-manage.sh; \
		./docker-manage.sh pull; \
	else \
		echo "⚠️ docker-manage.sh not found. Creating it..."; \
		echo "Please edit GITHUB_REPO variable in docker-manage.sh"; \
	fi

docker-run-ghcr:
	@echo "🚀 Running Docker container from GHCR..."
	@if [ -f ./docker-manage.sh ]; then \
		chmod +x ./docker-manage.sh; \
		./docker-manage.sh run; \
	else \
		echo "⚠️ docker-manage.sh not found"; \
	fi

docker-stop:
	@echo "⏹️ Stopping Docker container..."
	@if [ -f ./docker-manage.sh ]; then \
		chmod +x ./docker-manage.sh; \
		./docker-manage.sh stop; \
	else \
		docker stop tempfile || echo "Container not running"; \
	fi

docker-logs:
	@echo "📋 Showing Docker logs..."
	@if [ -f ./docker-manage.sh ]; then \
		chmod +x ./docker-manage.sh; \
		./docker-manage.sh logs; \
	else \
		docker logs -f tempfile || echo "Container not found"; \
	fi

docker-health:
	@echo "🏥 Checking Docker container health..."
	@if [ -f ./docker-manage.sh ]; then \
		chmod +x ./docker-manage.sh; \
		./docker-manage.sh health; \
	else \
		curl -f http://localhost:3000/health || echo "Health check failed"; \
	fi

docker-clean:
	@echo "🧹 Cleaning up Docker container..."
	@if [ -f ./docker-manage.sh ]; then \
		chmod +x ./docker-manage.sh; \
		./docker-manage.sh clean; \
	else \
		docker rm -f tempfile || echo "Container not found"; \
	fi

docker-update:
	@echo "⬆️ Updating to latest version..."
	@if [ -f ./docker-manage.sh ]; then \
		chmod +x ./docker-manage.sh; \
		./docker-manage.sh update; \
	else \
		echo "⚠️ docker-manage.sh not found"; \
	fi

docker-dev-ghcr:
	@echo "👨‍💻 Running development version from GHCR..."
	@if [ -f ./docker-manage.sh ]; then \
		chmod +x ./docker-manage.sh; \
		./docker-manage.sh dev; \
	else \
		echo "⚠️ docker-manage.sh not found"; \
	fi

docker-prod-ghcr:
	@echo "🏭 Running production version from GHCR..."
	@if [ -f ./docker-manage.sh ]; then \
		chmod +x ./docker-manage.sh; \
		./docker-manage.sh prod; \
	else \
		echo "⚠️ docker-manage.sh not found"; \
	fi

## Docker Compose targets
.PHONY: compose-up compose-down compose-logs compose-pull compose-restart

compose-up:
	@echo "🐳 Starting services with Docker Compose..."
	docker-compose up -d

compose-down:
	@echo "⏹️ Stopping services with Docker Compose..."
	docker-compose down

compose-logs:
	@echo "📋 Showing Docker Compose logs..."
	docker-compose logs -f

compose-pull:
	@echo "⬇️ Pulling latest images for compose..."
	docker-compose pull

compose-restart:
	@echo "🔄 Restarting services..."
	docker-compose restart

## Start Redis for development
redis-dev:
	@echo "🐳 Starting Redis for development..."
	@if ! docker ps | grep -q redis-dev; then \
		docker run -d --name redis-dev -p 6379:6379 redis:7-alpine; \
		echo "✅ Redis started on port 6379"; \
	else \
		echo "✅ Redis already running"; \
	fi

## Stop Redis development container
redis-stop:
	@echo "🐳 Stopping Redis development container..."
	@docker stop redis-dev && docker rm redis-dev || echo "Redis container not running"

## Test Docker health check
docker-health:
	@echo "🏥 Testing Docker health check..."
	curl -f http://localhost:3000/health || echo "Health check failed"

## Setup development environment
setup:
	@echo "🛠️ Setting up development environment..."
	@if [ -f ./dev.sh ]; then \
		chmod +x ./dev.sh; \
		./dev.sh setup; \
	else \
		echo "📦 Installing development tools..."; \
		go install github.com/cosmtrek/air@latest; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		echo "📋 Copying environment file..."; \
		cp .env.example .env; \
	fi
	@echo "✅ Development setup complete"

## Run with rate limiting enabled (memory store)
dev-ratelimit:
	@echo "🚀 Starting development server with rate limiting..."
	ENABLE_RATE_LIMIT=true $(GOCMD) run $(MAIN_PATH)

## Run with Redis rate limiting
dev-redis:
	@echo "🚀 Starting development server with Redis rate limiting..."
	@make redis-dev
	ENABLE_RATE_LIMIT=true RATE_LIMIT_STORE=redis $(GOCMD) run $(MAIN_PATH)

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
	@echo "🏗️  Development:"
	@echo "  setup        - Setup development environment"
	@echo "  dev          - Run in development mode"
	@echo "  dev-ratelimit - Run with rate limiting (memory store)"
	@echo "  dev-redis    - Run with Redis rate limiting"
	@echo ""
	@echo "🔨 Build & Run:"
	@echo "  build        - Build the application"
	@echo "  run          - Build and run the application"
	@echo "  build-all    - Build for all platforms"
	@echo "  ci-build     - Run full CI build (test + build-all)"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test         - Run all tests"
	@echo "  test-ratelimit - Run rate limiting tests only"
	@echo "  test-redis   - Run Redis integration tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo ""
	@echo "🔍 Quality:"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  security     - Run security scan"
	@echo "  check        - Run all quality checks (fmt, lint, security, test)"
	@echo ""
	@echo "🐳 Docker (Local Build):"
	@echo "  docker-build - Build Docker image locally"
	@echo "  docker-run   - Run local Docker container"
	@echo ""
	@echo "📶 Docker (GHCR):"
	@echo "  docker-pull      - Pull latest image from GHCR"
	@echo "  docker-run-ghcr  - Run container from GHCR"
	@echo "  docker-stop      - Stop running container"
	@echo "  docker-logs      - Show container logs"
	@echo "  docker-health    - Check container health"
	@echo "  docker-clean     - Remove container"
	@echo "  docker-update    - Update to latest version"
	@echo "  docker-dev-ghcr  - Run development version"
	@echo "  docker-prod-ghcr - Run production version"
	@echo ""
	@echo "📄 Docker Compose:"
	@echo "  compose-up       - Start services with compose"
	@echo "  compose-down     - Stop services"
	@echo "  compose-logs     - Show compose logs"
	@echo "  compose-pull     - Pull latest images"
	@echo "  compose-restart  - Restart services"
	@echo ""
	@echo "🗄️  Redis:"
	@echo "  redis-dev    - Start Redis for development"
	@echo "  redis-stop   - Stop Redis development container"
	@echo ""
	@echo "🛠️  Utilities:"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  commit       - Interactive commit helper (conventional commits)"
	@echo "  release-prep - Prepare for release (maintainers only)"
	@echo "  help         - Show this help"
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make setup            # First time setup"
	@echo "  make dev              # Start development"
	@echo "  make dev-ratelimit    # Test rate limiting"
	@echo "  make check            # Run quality checks"
	@echo "  make docker-run-ghcr  # Run from GHCR"
	@echo "  make compose-up       # Run with Docker Compose"
	@echo ""
