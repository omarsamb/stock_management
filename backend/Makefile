# Database Migration Commands

# Variables
MIGRATIONS_DIR=db/migrations
DATABASE_URL=postgres://postgres:postgres@localhost:5432/stock_management?sslmode=disable

# Install golang-migrate if not present
.PHONY: install-migrate
install-migrate:
	@which migrate > /dev/null || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create a new migration
.PHONY: migrate-create
migrate-create: install-migrate
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

# Run migrations up
.PHONY: migrate-up
migrate-up: install-migrate
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

# Run migrations down
.PHONY: migrate-down
migrate-down: install-migrate
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down

# Force migration version (use with caution)
.PHONY: migrate-force
migrate-force: install-migrate
	@read -p "Enter version to force: " version; \
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $$version

# Check migration version
.PHONY: migrate-version
migrate-version: install-migrate
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make install-migrate  - Install golang-migrate tool"
	@echo "  make migrate-create   - Create a new migration file"
	@echo "  make migrate-up       - Run all pending migrations"
	@echo "  make migrate-down     - Rollback last migration"
	@echo "  make migrate-version  - Show current migration version"
	@echo "  make migrate-force    - Force migration to specific version"

# Build the backend
.PHONY: build
build:
	go build -o backend

# Run the backend locally
.PHONY: run
run: build
	./backend

# Default target
.DEFAULT_GOAL := help