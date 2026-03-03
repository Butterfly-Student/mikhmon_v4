#!/bin/bash

# Script untuk menjalankan GoMikhmon dengan Podman tanpa compose

set -e

NETWORK_NAME="mikhmon"
POSTGRES_VOLUME="postgres_data"

# Fungsi untuk membuat network
create_network() {
    if ! podman network exists "$NETWORK_NAME"; then
        echo "Creating network: $NETWORK_NAME"
        podman network create "$NETWORK_NAME"
    else
        echo "Network $NETWORK_NAME already exists"
    fi
}

# Fungsi untuk membuat volume
create_volume() {
    if ! podman volume exists "$POSTGRES_VOLUME"; then
        echo "Creating volume: $POSTGRES_VOLUME"
        podman volume create "$POSTGRES_VOLUME"
    else
        echo "Volume $POSTGRES_VOLUME already exists"
    fi
}

# Start services
start() {
    echo "Starting GoMikhmon services..."
    create_network
    create_volume

    # Start PostgreSQL
    if ! podman ps --format "{{.Names}}" | grep -q "^gomikhmon-postgres$"; then
        if podman ps -a --format "{{.Names}}" | grep -q "^gomikhmon-postgres$"; then
            echo "Starting existing postgres container..."
            podman start gomikhmon-postgres
        else
            echo "Creating and starting postgres container..."
            podman run -d \
                --name gomikhmon-postgres \
                --network "$NETWORK_NAME" \
                -p 5432:5432 \
                -e POSTGRES_USER=mikhmon \
                -e POSTGRES_PASSWORD=mikhmon \
                -e POSTGRES_DB=mikhmon \
                -v "$POSTGRES_VOLUME:/var/lib/postgresql/data" \
                --restart unless-stopped \
                docker.io/library/postgres:15-alpine
        fi
    else
        echo "Postgres already running"
    fi

    # Tunggu postgres siap
    echo "Waiting for postgres to be ready..."
    sleep 3
    podman exec gomikhmon-postgres pg_isready -U mikhmon || sleep 5

    # Start Redis
    if ! podman ps --format "{{.Names}}" | grep -q "^gomikhmon-redis$"; then
        if podman ps -a --format "{{.Names}}" | grep -q "^gomikhmon-redis$"; then
            echo "Starting existing redis container..."
            podman start gomikhmon-redis
        else
            echo "Creating and starting redis container..."
            podman run -d \
                --name gomikhmon-redis \
                --network "$NETWORK_NAME" \
                -p 6379:6379 \
                --restart unless-stopped \
                docker.io/library/redis:7-alpine
        fi
    else
        echo "Redis already running"
    fi

    echo ""
    echo "Services started successfully!"
    echo "PostgreSQL: localhost:5432"
    echo "Redis:      localhost:6379"
}

# Stop services
stop() {
    echo "Stopping GoMikhmon services..."
    podman stop gomikhmon-redis 2>/dev/null || true
    podman stop gomikhmon-postgres 2>/dev/null || true
    echo "Services stopped"
}

# Restart services
restart() {
    stop
    sleep 2
    start
}

# Show status
status() {
    echo "Container Status:"
    podman ps -a --filter "name=gomikhmon" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
}

# View logs
logs() {
    if [ -z "$1" ]; then
        echo "Usage: $0 logs [postgres|redis]"
        exit 1
    fi
    podman logs -f "gomikhmon-$1"
}

# Remove containers and volumes
down() {
    echo "Removing containers..."
    podman rm -f gomikhmon-redis 2>/dev/null || true
    podman rm -f gomikhmon-postgres 2>/dev/null || true
    
    echo "Removing network..."
    podman network rm "$NETWORK_NAME" 2>/dev/null || true
    
    echo "Note: Volume '$POSTGRES_VOLUME' preserved. Use 'clean' to remove all data."
}

# Clean everything including volumes
clean() {
    down
    echo "Removing volume..."
    podman volume rm "$POSTGRES_VOLUME" 2>/dev/null || true
    echo "All data cleaned"
}

# Help
help() {
    echo "Usage: $0 {start|stop|restart|status|logs|down|clean}"
    echo ""
    echo "Commands:"
    echo "  start    - Start all services"
    echo "  stop     - Stop all services"
    echo "  restart  - Restart all services"
    echo "  status   - Show container status"
    echo "  logs     - View logs (e.g., $0 logs postgres)"
    echo "  down     - Stop and remove containers"
    echo "  clean    - Stop, remove containers and delete all data"
}

# Main
case "${1:-}" in
    start|up)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    status|ps)
        status
        ;;
    logs)
        logs "$2"
        ;;
    down)
        down
        ;;
    clean)
        clean
        ;;
    *)
        help
        ;;
esac
