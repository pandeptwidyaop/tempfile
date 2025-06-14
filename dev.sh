#!/bin/bash

# 🚀 TempFile Development Helper Script
# This script helps with common development tasks

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go 1.24 or higher."
        exit 1
    fi
    
    GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | head -1)
    log_info "Found Go version: $GO_VERSION"
}

# Setup development environment
setup() {
    log_info "Setting up development environment..."
    
    check_go
    
    # Download dependencies
    log_info "Downloading Go dependencies..."
    go mod download
    go mod verify
    
    # Create necessary directories
    mkdir -p uploads
    
    # Copy environment file if it doesn't exist
    if [ ! -f .env ]; then
        cp .env.example .env
        log_success "Created .env file from .env.example"
        log_warning "Please review and customize your .env file"
    fi
    
    # Install development tools (optional)
    read -p "Install development tools (air, golangci-lint, gosec)? [y/N]: " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        log_info "Installing development tools..."
        go install github.com/cosmtrek/air@latest
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
        log_success "Development tools installed"
    fi
    
    log_success "Development environment setup complete!"
}

# Run tests
test() {
    log_info "Running tests..."
    go test -v ./...
    log_success "All tests passed!"
}

# Run tests with coverage
test_coverage() {
    log_info "Running tests with coverage..."
    go test -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    log_success "Coverage report generated: coverage.html"
}

# Format code
format() {
    log_info "Formatting Go code..."
    go fmt ./...
    log_success "Code formatted!"
}

# Lint code
lint() {
    if command -v golangci-lint &> /dev/null; then
        log_info "Running golangci-lint..."
        golangci-lint run
        log_success "Linting complete!"
    else
        log_warning "golangci-lint not installed. Installing..."
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        log_info "Running golangci-lint..."
        golangci-lint run
        log_success "Linting complete!"
    fi
}

# Security scan
security() {
    if command -v gosec &> /dev/null; then
        log_info "Running security scan..."
        gosec ./...
        log_success "Security scan complete!"
    else
        log_warning "gosec not installed. Installing..."
        go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
        log_info "Running security scan..."
        gosec ./...
        log_success "Security scan complete!"
    fi
}

# Build application
build() {
    log_info "Building application..."
    
    # Create bin directory if it doesn't exist
    mkdir -p bin
    
    # Build with version info
    VERSION=${1:-"dev"}
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
    
    LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}"
    
    go build -ldflags="${LDFLAGS}" -o bin/tempfile ./cmd/server
    
    log_success "Build complete: bin/tempfile"
    log_info "Version: $VERSION, Commit: $COMMIT"
}

# Build for multiple platforms
build_all() {
    log_info "Building for multiple platforms..."
    
    mkdir -p dist
    
    VERSION=${1:-"dev"}
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S_UTC')
    
    LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}"
    
    # Build for different platforms
    platforms=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")
    
    for platform in "${platforms[@]}"; do
        platform_split=(${platform//\// })
        GOOS=${platform_split[0]}
        GOARCH=${platform_split[1]}
        
        output_name="tempfile-${GOOS}-${GOARCH}"
        if [ $GOOS = "windows" ]; then
            output_name+='.exe'
        fi
        
        log_info "Building for $GOOS/$GOARCH..."
        env GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/$output_name ./cmd/server
    done
    
    log_success "Multi-platform build complete! Check dist/ directory"
}

# Run application in development mode
dev() {
    if command -v air &> /dev/null; then
        log_info "Starting development server with live reload..."
        air
    else
        log_info "Starting development server..."
        go run cmd/server/main.go
    fi
}

# Clean build artifacts
clean() {
    log_info "Cleaning build artifacts..."
    rm -rf bin/ dist/ coverage.out coverage.html
    log_success "Clean complete!"
}

# Run all quality checks
check() {
    log_info "Running all quality checks..."
    format
    lint
    security  
    test
    log_success "All checks passed! ✨"
}

# Commit helper with conventional commits
commit() {
    log_info "Interactive commit helper..."
    
    # Check if there are staged changes
    if ! git diff --cached --quiet; then
        log_error "No staged changes found. Please stage your changes first:"
        echo "  git add <files>"
        exit 1
    fi
    
    echo "Select commit type:"
    echo "1) feat     - New feature"
    echo "2) fix      - Bug fix"
    echo "3) docs     - Documentation"
    echo "4) style    - Code style"
    echo "5) refactor - Code refactoring"
    echo "6) test     - Tests"
    echo "7) chore    - Maintenance"
    echo "8) ci       - CI/CD"
    
    read -p "Enter choice [1-8]: " choice
    
    case $choice in
        1) type="feat";;
        2) type="fix";;
        3) type="docs";;
        4) type="style";;
        5) type="refactor";;
        6) type="test";;
        7) type="chore";;
        8) type="ci";;
        *) log_error "Invalid choice"; exit 1;;
    esac
    
    read -p "Enter scope (optional, e.g., api, ui): " scope
    read -p "Enter description: " description
    read -p "Is this a breaking change? [y/N]: " breaking
    
    if [ -n "$scope" ]; then
        scope="($scope)"
    fi
    
    if [[ $breaking =~ ^[Yy]$ ]]; then
        type="${type}!"
    fi
    
    commit_msg="${type}${scope}: ${description}"
    
    echo
    log_info "Commit message: $commit_msg"
    read -p "Proceed with commit? [Y/n]: " proceed
    
    if [[ ! $proceed =~ ^[Nn]$ ]]; then
        git commit -m "$commit_msg"
        log_success "Commit created successfully!"
    else
        log_info "Commit cancelled"
    fi
}

# Release helper (for maintainers)
release() {
    log_info "Preparing release..."
    
    # Check if we're on main branch
    current_branch=$(git branch --show-current)
    if [ "$current_branch" != "main" ]; then
        log_error "Please switch to main branch before releasing"
        exit 1
    fi
    
    # Check if working directory is clean
    if ! git diff --quiet || ! git diff --cached --quiet; then
        log_error "Working directory is not clean. Please commit or stash changes."
        exit 1
    fi
    
    # Pull latest changes
    log_info "Pulling latest changes..."
    git pull origin main
    
    # Run quality checks
    check
    
    log_success "Ready to release! Push to main branch to trigger automatic release."
    log_info "GitHub Actions will handle versioning and release creation."
}

# Show help
show_help() {
    echo "🚀 TempFile Development Helper"
    echo
    echo "Usage: $0 <command> [options]"
    echo
    echo "Commands:"
    echo "  setup           Setup development environment"
    echo "  dev             Start development server"
    echo "  test            Run tests"
    echo "  test-coverage   Run tests with coverage report"
    echo "  format          Format Go code"
    echo "  lint            Run linter"
    echo "  security        Run security scan"
    echo "  build [version] Build application (default: dev)"
    echo "  build-all [ver] Build for all platforms"
    echo "  clean           Clean build artifacts"
    echo "  check           Run all quality checks"
    echo "  commit          Interactive commit helper"
    echo "  release         Prepare release (maintainers)"
    echo "  help            Show this help"
    echo
    echo "Examples:"
    echo "  $0 setup           # Setup development environment"
    echo "  $0 dev             # Start development server"
    echo "  $0 build v1.0.0    # Build with version v1.0.0"
    echo "  $0 check           # Run all quality checks"
}

# Main script logic
case "${1:-help}" in
    setup)
        setup
        ;;
    dev)
        dev
        ;;
    test)
        test
        ;;
    test-coverage)
        test_coverage
        ;;
    format)
        format
        ;;
    lint)
        lint
        ;;
    security)
        security
        ;;
    build)
        build "$2"
        ;;
    build-all)
        build_all "$2"
        ;;
    clean)
        clean
        ;;
    check)
        check
        ;;
    commit)
        commit
        ;;
    release)
        release
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        log_error "Unknown command: $1"
        echo
        show_help
        exit 1
        ;;
esac
