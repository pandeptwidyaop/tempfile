#!/bin/bash

# Tempfile Docker Management Script
# This script helps manage Tempfile Docker containers

set -e

# Configuration
GITHUB_REPO="$(git remote get-url origin | sed 's/.*github.com[\/:]//g' | sed 's/\.git$//')"
IMAGE_NAME="ghcr.io/${GITHUB_REPO,,}"
CONTAINER_NAME="tempfile"
DEFAULT_PORT="3000"
UPLOADS_DIR="./uploads"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Show usage
show_usage() {
    cat << EOF
🐳 Tempfile Docker Management Script

Usage: $0 [COMMAND] [OPTIONS]

Commands:
  pull [TAG]           Pull Docker image (default: latest)
  run [TAG] [PORT]     Run container (default: latest, 3000)
  stop                 Stop running container
  restart              Restart container
  logs                 Show container logs
  shell                Access container shell
  health               Check container health
  clean                Remove stopped container
  update               Update to latest version
  dev                  Run development version
  prod                 Run production version with optimizations

Examples:
  $0 pull                    # Pull latest image
  $0 run v1.2.0 8080        # Run specific version on port 8080
  $0 dev                     # Run development version
  $0 prod                    # Run production version
  $0 logs                    # Show logs
  $0 health                  # Check health

Environment Variables:
  TEMPFILE_PORT     Port to run on (default: 3000)
  TEMPFILE_ENV      Environment mode (development/production)
  TEMPFILE_DEBUG    Enable debug mode (true/false)
  PUBLIC_URL        Public URL for the application

EOF
}

# Check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
}

# Check if container exists
container_exists() {
    docker ps -a --format "table {{.Names}}" | grep -q "^${CONTAINER_NAME}$"
}

# Check if container is running
container_running() {
    docker ps --format "table {{.Names}}" | grep -q "^${CONTAINER_NAME}$"
}

# Pull Docker image
pull_image() {
    local tag=${1:-latest}
    log_info "Pulling Docker image: ${IMAGE_NAME}:${tag}"
    
    if docker pull "${IMAGE_NAME}:${tag}"; then
        log_success "Successfully pulled ${IMAGE_NAME}:${tag}"
    else
        log_error "Failed to pull image"
        exit 1
    fi
}

# Run container
run_container() {
    local tag=${1:-latest}
    local port=${2:-${TEMPFILE_PORT:-$DEFAULT_PORT}}
    
    # Stop existing container if running
    if container_running; then
        log_warning "Container is already running. Stopping it first..."
        stop_container
    fi
    
    # Remove existing container if exists
    if container_exists; then
        log_info "Removing existing container..."
        docker rm "${CONTAINER_NAME}" > /dev/null
    fi
    
    # Create uploads directory
    mkdir -p "${UPLOADS_DIR}"
    
    # Set environment variables
    local env_vars=""
    if [[ -n "${TEMPFILE_ENV}" ]]; then
        env_vars="${env_vars} -e APP_ENV=${TEMPFILE_ENV}"
    fi
    if [[ -n "${TEMPFILE_DEBUG}" ]]; then
        env_vars="${env_vars} -e DEBUG=${TEMPFILE_DEBUG}"
    fi
    if [[ -n "${PUBLIC_URL}" ]]; then
        env_vars="${env_vars} -e PUBLIC_URL=${PUBLIC_URL}"
    fi
    
    log_info "Running container: ${IMAGE_NAME}:${tag} on port ${port}"
    
    if docker run -d \
        --name "${CONTAINER_NAME}" \
        --restart unless-stopped \
        -p "${port}:3000" \
        -v "$(pwd)/${UPLOADS_DIR}:/app/uploads" \
        ${env_vars} \
        "${IMAGE_NAME}:${tag}"; then
        
        log_success "Container started successfully"
        log_info "Application will be available at: http://localhost:${port}"
        log_info "Waiting for container to be ready..."
        
        # Wait for container to be healthy
        local max_attempts=30
        local attempt=1
        while [[ $attempt -le $max_attempts ]]; do
            if curl -s "http://localhost:${port}/health" > /dev/null 2>&1; then
                log_success "Container is healthy and ready!"
                break
            fi
            
            if [[ $attempt -eq $max_attempts ]]; then
                log_warning "Container may not be fully ready yet. Check logs with: $0 logs"
                break
            fi
            
            echo -n "."
            sleep 2
            ((attempt++))
        done
        echo
    else
        log_error "Failed to start container"
        exit 1
    fi
}

# Stop container
stop_container() {
    if container_running; then
        log_info "Stopping container..."
        if docker stop "${CONTAINER_NAME}"; then
            log_success "Container stopped"
        else
            log_error "Failed to stop container"
            exit 1
        fi
    else
        log_warning "Container is not running"
    fi
}

# Restart container
restart_container() {
    if container_exists; then
        log_info "Restarting container..."
        if docker restart "${CONTAINER_NAME}"; then
            log_success "Container restarted"
        else
            log_error "Failed to restart container"
            exit 1
        fi
    else
        log_error "Container does not exist. Use 'run' command first."
        exit 1
    fi
}

# Show logs
show_logs() {
    if container_exists; then
        log_info "Showing container logs..."
        docker logs -f "${CONTAINER_NAME}"
    else
        log_error "Container does not exist"
        exit 1
    fi
}

# Access shell
access_shell() {
    if container_running; then
        log_info "Accessing container shell..."
        docker exec -it "${CONTAINER_NAME}" sh
    else
        log_error "Container is not running"
        exit 1
    fi
}

# Check health
check_health() {
    if container_running; then
        local port=$(docker port "${CONTAINER_NAME}" 3000 | cut -d: -f2)
        log_info "Checking container health..."
        
        if curl -s "http://localhost:${port}/health" | grep -q "ok"; then
            log_success "Container is healthy"
            
            # Show container info
            echo
            log_info "Container Information:"
            docker ps --filter "name=${CONTAINER_NAME}" --format "table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}"
            
            # Show resource usage
            echo
            log_info "Resource Usage:"
            docker stats "${CONTAINER_NAME}" --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}"
        else
            log_error "Container is not healthy"
            exit 1
        fi
    else
        log_error "Container is not running"
        exit 1
    fi
}

# Clean up
clean_container() {
    if container_exists; then
        if container_running; then
            log_info "Stopping container first..."
            stop_container
        fi
        
        log_info "Removing container..."
        if docker rm "${CONTAINER_NAME}"; then
            log_success "Container removed"
        else
            log_error "Failed to remove container"
            exit 1
        fi
    else
        log_warning "Container does not exist"
    fi
}

# Update to latest version
update_container() {
    log_info "Updating to latest version..."
    
    # Pull latest image
    pull_image "latest"
    
    # Restart with new image
    if container_running; then
        local port=$(docker port "${CONTAINER_NAME}" 3000 | cut -d: -f2)
        stop_container
        clean_container
        run_container "latest" "${port}"
    else
        log_info "Container is not running. Use 'run' command to start with latest image."
    fi
}

# Run development version
run_dev() {
    export TEMPFILE_ENV="development"
    export TEMPFILE_DEBUG="true"
    log_info "Running development version..."
    pull_image "develop"
    run_container "develop" "${1:-3000}"
}

# Run production version
run_prod() {
    export TEMPFILE_ENV="production"
    export TEMPFILE_DEBUG="false"
    log_info "Running production version with optimizations..."
    pull_image "latest"
    
    local port=${1:-${TEMPFILE_PORT:-$DEFAULT_PORT}}
    
    # Stop existing container if running
    if container_running; then
        stop_container
    fi
    
    # Remove existing container if exists
    if container_exists; then
        docker rm "${CONTAINER_NAME}" > /dev/null
    fi
    
    # Create uploads directory
    mkdir -p "${UPLOADS_DIR}"
    
    log_info "Running optimized production container..."
    
    if docker run -d \
        --name "${CONTAINER_NAME}" \
        --restart unless-stopped \
        --memory="512m" \
        --cpus="1.0" \
        --security-opt no-new-privileges:true \
        -p "${port}:3000" \
        -v "$(pwd)/${UPLOADS_DIR}:/app/uploads" \
        -e APP_ENV=production \
        -e DEBUG=false \
        ${PUBLIC_URL:+-e PUBLIC_URL="${PUBLIC_URL}"} \
        "${IMAGE_NAME}:latest"; then
        
        log_success "Production container started successfully"
        log_info "Application available at: http://localhost:${port}"
    else
        log_error "Failed to start production container"
        exit 1
    fi
}

# Main script logic
main() {
    check_docker
    
    case "${1:-}" in
        "pull")
            pull_image "$2"
            ;;
        "run")
            pull_image "$2"
            run_container "$2" "$3"
            ;;
        "stop")
            stop_container
            ;;
        "restart")
            restart_container
            ;;
        "logs")
            show_logs
            ;;
        "shell")
            access_shell
            ;;
        "health")
            check_health
            ;;
        "clean")
            clean_container
            ;;
        "update")
            update_container
            ;;
        "dev")
            run_dev "$2"
            ;;
        "prod")
            run_prod "$2"
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        "")
            log_error "No command specified"
            show_usage
            exit 1
            ;;
        *)
            log_error "Unknown command: $1"
            show_usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
