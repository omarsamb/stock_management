#!/bin/bash

# Database Migration Script for Stock Management System
# This script runs database migrations using golang-migrate

set -e

# Configuration
MIGRATIONS_DIR="db/migrations"
DB_HOST="${DB_HOST:-localhost}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-stock_management}"
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=disable"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if migrate command is available
if ! command -v migrate &> /dev/null; then
    log_error "golang-migrate is not installed. Please install it first:"
    echo "go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Show current migration version
log_info "Checking current migration version..."
migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version || log_warn "No migrations applied yet"

# Run migrations
log_info "Running database migrations..."
migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up

log_info "Migrations completed successfully!"

# Show final migration version
log_info "Current migration version:"
migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version