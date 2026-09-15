#!/bin/bash
# Octopus Development Environment Startup Script
# One-shot script to start all services for local development

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_docker() {
    if command -v docker &> /dev/null && docker info &> /dev/null; then
        return 0
    fi
    return 1
}

check_local_services() {
    local mysql_ok=false
    local redis_ok=false
    
    if command -v mysqladmin &> /dev/null && mysqladmin ping -h localhost -u root -poctopus_root_pwd &> /dev/null 2>&1; then
        mysql_ok=true
    fi
    
    if command -v redis-cli &> /dev/null && redis-cli ping &> /dev/null 2>&1; then
        redis_ok=true
    fi
    
    if $mysql_ok && $redis_ok; then
        return 0
    fi
    return 1
}

start_infrastructure_docker() {
    log_info "Starting infrastructure via Docker Compose..."
    cd "$PROJECT_ROOT"
    docker compose up -d
    
    log_info "Waiting for MySQL to be ready..."
    local max_attempts=30
    local attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if docker exec octopus-mysql mysqladmin ping -h localhost -u root -poctopus_root_pwd &> /dev/null 2>&1; then
            log_success "MySQL is ready"
            break
        fi
        attempt=$((attempt + 1))
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        log_error "MySQL failed to start within timeout"
        exit 1
    fi
    
    log_info "Waiting for Redis to be ready..."
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if docker exec octopus-redis redis-cli ping &> /dev/null 2>&1; then
            log_success "Redis is ready"
            break
        fi
        attempt=$((attempt + 1))
        sleep 1
    done
    
    if [ $attempt -eq $max_attempts ]; then
        log_error "Redis failed to start within timeout"
        exit 1
    fi
}

start_infrastructure_local() {
    log_info "Starting local MySQL and Redis services..."
    
    if command -v service &> /dev/null; then
        sudo service mysql start 2>/dev/null || true
        sudo service redis-server start 2>/dev/null || true
    elif command -v systemctl &> /dev/null; then
        sudo systemctl start mysql 2>/dev/null || sudo systemctl start mysqld 2>/dev/null || true
        sudo systemctl start redis 2>/dev/null || sudo systemctl start redis-server 2>/dev/null || true
    fi
    
    sleep 2
    
    if ! check_local_services; then
        log_error "Failed to start local MySQL/Redis services"
        log_info "Please install MySQL and Redis manually or use Docker"
        exit 1
    fi
    
    log_success "Local MySQL and Redis are running"
}

init_database() {
    log_info "Initializing database schema..."
    
    local mysql_cmd="mysql -h localhost -u root -poctopus_root_pwd"
    
    if check_docker && docker ps | grep -q octopus-mysql; then
        mysql_cmd="docker exec -i octopus-mysql mysql -u root -poctopus_root_pwd"
    fi
    
    $mysql_cmd -e "CREATE DATABASE IF NOT EXISTS octopus;" 2>/dev/null || true
    $mysql_cmd octopus < "$PROJECT_ROOT/docker/mysql/init/00_schema.sql" 2>/dev/null
    
    log_success "Database schema initialized"
}

start_backend() {
    log_info "Starting Go backend server..."
    cd "$PROJECT_ROOT/backend"
    
    go mod tidy
    
    if [ "$DETACH" = "true" ]; then
        nohup go run workflow.go -f etc/workflow-api.yaml > /tmp/octopus-backend.log 2>&1 &
        echo $! > /tmp/octopus-backend.pid
        sleep 3
        
        if curl -s http://localhost:8888/api/v1/health | grep -q "ok"; then
            log_success "Backend started (PID: $(cat /tmp/octopus-backend.pid))"
            log_info "Backend logs: /tmp/octopus-backend.log"
        else
            log_error "Backend failed to start. Check /tmp/octopus-backend.log"
            exit 1
        fi
    else
        log_info "Starting backend in foreground (use DETACH=true for background)..."
        go run workflow.go -f etc/workflow-api.yaml
    fi
}

start_worker() {
    log_info "Starting Python worker..."
    cd "$PROJECT_ROOT/worker"
    
    if [ ! -d "venv" ]; then
        python3 -m venv venv
    fi
    
    source venv/bin/activate
    pip install -q -r requirements.txt
    
    if [ "$DETACH" = "true" ]; then
        nohup python3 main.py > /tmp/octopus-worker.log 2>&1 &
        echo $! > /tmp/octopus-worker.pid
        sleep 2
        
        if ps -p $(cat /tmp/octopus-worker.pid) > /dev/null 2>&1; then
            log_success "Worker started (PID: $(cat /tmp/octopus-worker.pid))"
            log_info "Worker logs: /tmp/octopus-worker.log"
        else
            log_error "Worker failed to start. Check /tmp/octopus-worker.log"
            exit 1
        fi
    else
        log_info "Starting worker in foreground (use DETACH=true for background)..."
        python3 main.py
    fi
}

show_status() {
    echo ""
    echo "=========================================="
    echo "  Octopus Development Environment"
    echo "=========================================="
    echo ""
    echo "Services:"
    echo "  - MySQL:   localhost:3306"
    echo "  - Redis:   localhost:6379"
    echo "  - Backend: http://localhost:8888"
    echo ""
    echo "Health Check:"
    echo "  curl http://localhost:8888/api/v1/health"
    echo ""
    echo "Quick Test:"
    echo "  # Create workflow"
    echo "  curl -X POST http://localhost:8888/api/v1/workflows \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"name\":\"Test\",\"nodes\":[{\"id\":\"n1\",\"type\":\"script\",\"name\":\"N1\",\"position\":{\"x\":0,\"y\":0}}],\"entry_node_id\":\"n1\"}'"
    echo ""
    echo "  # Enable and run"
    echo "  curl -X POST http://localhost:8888/api/v1/workflows/1/enable"
    echo "  curl -X POST http://localhost:8888/api/v1/workflows/1/runs"
    echo ""
}

usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  infra      Start only MySQL and Redis (Docker or local)"
    echo "  backend    Start only the Go backend (requires infra)"
    echo "  worker     Start only the Python worker (requires infra + backend)"
    echo "  all        Start everything (default)"
    echo "  status     Show service status and helpful commands"
    echo ""
    echo "Environment Variables:"
    echo "  DETACH=true    Run services in background (default: foreground)"
    echo "  USE_DOCKER=1   Force Docker for infrastructure"
    echo ""
    echo "Examples:"
    echo "  $0                    # Start all services in foreground"
    echo "  DETACH=true $0 all    # Start all services in background"
    echo "  $0 infra              # Start only MySQL and Redis"
    echo ""
}

main() {
    local cmd="${1:-all}"
    
    case "$cmd" in
        infra)
            if [ "$USE_DOCKER" = "1" ] || check_docker; then
                start_infrastructure_docker
            else
                start_infrastructure_local
            fi
            init_database
            ;;
        backend)
            start_backend
            ;;
        worker)
            start_worker
            ;;
        all)
            if [ "$USE_DOCKER" = "1" ] || check_docker; then
                start_infrastructure_docker
            elif check_local_services; then
                log_success "Using existing local MySQL and Redis"
            else
                start_infrastructure_local
            fi
            init_database
            
            if [ "$DETACH" = "true" ]; then
                start_backend
                start_worker
                show_status
            else
                log_info "Starting backend (press Ctrl+C to stop)..."
                log_warn "To run worker, open another terminal and run: $0 worker"
                start_backend
            fi
            ;;
        status)
            show_status
            ;;
        -h|--help|help)
            usage
            ;;
        *)
            log_error "Unknown command: $cmd"
            usage
            exit 1
            ;;
    esac
}

main "$@"
