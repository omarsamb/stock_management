# Database Schema & Migrations

This document describes the database schema and migration setup for the stock management system.

## Overview

The system uses PostgreSQL with a multi-tenant architecture where all data is partitioned by `tenant_id`. Each tenant represents a separate business/shop using the system.

## Schema Design

### Core Tables

#### 1. Tenants
Represents individual businesses/shops using the system.
```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
```

#### 2. Users
Users associated with tenants, with role-based access control.
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin','manager','staff')),
    created_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
```

#### 3. Items
Inventory items managed by each tenant.
```sql
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    name TEXT NOT NULL,
    sku TEXT UNIQUE,
    quantity INT NOT NULL DEFAULT 0,
    min_quantity INT NOT NULL DEFAULT 0, -- for low-stock alerts
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
```

#### 4. Stock Movements
History log of all stock changes (increases/decreases).
```sql
CREATE TABLE stock_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    item_id UUID NOT NULL,
    change INT NOT NULL, -- positive or negative
    reason TEXT,         -- e.g., "sale", "restock"
    created_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (item_id) REFERENCES items(id)
);
```

#### 5. Supplier Orders
Orders placed with suppliers.
```sql
CREATE TABLE supplier_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    supplier_name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending','shipped','received','cancelled')),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
```

#### 6. Supplier Order Items
Items within supplier orders.
```sql
CREATE TABLE supplier_order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL,
    item_id UUID NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES supplier_orders(id),
    FOREIGN KEY (item_id) REFERENCES items(id)
);
```

## Migration System

We use [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations.

### Migration Files

Migrations are stored in `db/migrations/` with the following naming convention:
- `{version}_{description}.up.sql` - Apply migration
- `{version}_{description}.down.sql` - Rollback migration

Current migrations:
1. `000001_create_initial_schema` - Creates tenants and users tables
2. `000002_create_stock_management_tables` - Creates items, stock_movements, supplier_orders, supplier_order_items tables

### Running Migrations

#### Using Make (Recommended)
```bash
# Install golang-migrate (if not already installed)
make install-migrate

# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check current migration version
make migrate-version

# Create new migration
make migrate-create
```

#### Using Docker Compose
Migrations run automatically when starting the application:
```bash
docker compose up
```

The `migrate` service runs first and applies all pending migrations before the backend starts.

#### Manual Migration
```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
./migrate.sh

# Or use migrate directly
migrate -path db/migrations -database "postgres://user:pass@host:port/dbname?sslmode=disable" up
```

### Migration Environment Variables

- `DB_HOST` - Database host (default: localhost)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: stock_management)

## Indexes

The following indexes are created for optimal performance:

**Users Table:**
- `idx_users_tenant_id` - For tenant-based queries
- `idx_users_email` - For login lookups

**Items Table:**
- `idx_items_tenant_id` - For tenant-based queries
- `idx_items_sku` - For SKU lookups

**Stock Movements Table:**
- `idx_stock_movements_tenant_id` - For tenant-based queries
- `idx_stock_movements_item_id` - For item history queries

**Supplier Orders Table:**
- `idx_supplier_orders_tenant_id` - For tenant-based queries

**Supplier Order Items Table:**
- `idx_supplier_order_items_order_id` - For order detail queries
- `idx_supplier_order_items_item_id` - For item usage queries

## Multi-Tenancy

All data tables include a `tenant_id` column with foreign key constraints to the `tenants` table. This ensures:

1. **Data Isolation**: Each tenant can only access their own data
2. **Referential Integrity**: All data belongs to a valid tenant
3. **Performance**: Indexes on `tenant_id` enable efficient tenant-scoped queries

## Security Considerations

1. **UUID Primary Keys**: Prevent enumeration attacks
2. **Role-Based Access**: Users have specific roles (admin, manager, staff)
3. **Foreign Key Constraints**: Maintain data integrity
4. **Unique Email Constraint**: Prevent duplicate user accounts
5. **Password Hashing**: Passwords are never stored in plain text