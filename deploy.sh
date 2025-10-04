#!/bin/bash

# Build and deployment script for allchat-BE

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status()   { echo -e "${GREEN}[INFO]${NC} $1"; }
print_warning()  { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error()    { echo -e "${RED}[ERROR]${NC} $1"; }

# ---- Prerequisites ----
if ! command -v docker &>/dev/null; then
  print_error "Docker is not installed. Please install Docker first."
  exit 1
fi

# Detect Compose command (v2 plugin preferred, then v1 binary)
if docker compose version &>/dev/null; then
  COMPOSE="docker compose"
elif command -v docker-compose &>/dev/null; then
  COMPOSE="docker-compose"
else
  print_error "Docker Compose is not installed. Install either 'docker compose' (v2) or 'docker-compose' (v1)."
  exit 1
fi

# ---- Paths ----
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "$SCRIPT_DIR"

COMPOSE_FILE="docker/docker-compose.yml"
ENV_FILE=".env"

if [ ! -f "$COMPOSE_FILE" ]; then
  print_error "Compose file not found at: $COMPOSE_FILE"
  exit 1
fi

if [ ! -f "$ENV_FILE" ]; then
  print_warning ".env not found at: $ENV_FILE"
  print_warning "Create an .env file (you can copy from .env.example if available)."
  exit 1
fi

# ---- Actions ----
build() {
  print_status "Building Docker images..."
  $COMPOSE -f "$COMPOSE_FILE" --env-file "$ENV_FILE" build --no-cache
  print_status "Build completed successfully!"
}

start() {
  print_status "Starting services..."
  $COMPOSE -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d
  print_status "Services started successfully!"
}

stop() {
  print_status "Stopping services..."
  $COMPOSE -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down
  print_status "Services stopped successfully!"
}

restart() {
  stop
  start
}

logs() {
  local service_name="${2:-}"
  if [ -n "$service_name" ]; then
    print_status "Showing logs for service: $service_name"
    $COMPOSE -f "$COMPOSE_FILE" --env-file "$ENV_FILE" logs -f "$service_name"
  else
    print_status "Showing logs for all services"
    $COMPOSE -f "$COMPOSE_FILE" --env-file "$ENV_FILE" logs -f
  fi
}

status() {
  $COMPOSE -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps
}

clean() {
  print_status "Cleaning up Docker resources..."
  $COMPOSE -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down -v --rmi all --remove-orphans
  docker system prune -f
  print_status "Cleanup completed!"
}

help() {
  cat <<EOF
Usage: $0 {build|start|stop|restart|logs|status|clean|help}

Commands:
  build     - Build Docker images
  start     - Start all services
  stop      - Stop all services
  restart   - Restart all services
  logs      - Show logs from all services
  logs <service> - Show logs from specific service (e.g., logs app, logs postgres)
  status    - Show status of all services
  clean     - Clean up all Docker resources
  help      - Show this help message


Detected compose command: $COMPOSE
Compose file: $COMPOSE_FILE
Env file: $ENV_FILE
EOF
}

case "${1:-help}" in
  build)   build ;;
  start)   start ;;
  stop)    stop ;;
  restart) restart ;;
  logs)    logs "$@" ;;
  status)  status ;;
  clean)   clean ;;
  help|*)  help; exit 0 ;;
esac