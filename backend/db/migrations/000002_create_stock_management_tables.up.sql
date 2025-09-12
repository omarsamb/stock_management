-- Items table with tenant_id for multi-tenancy
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

-- Stock movements history log
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

-- Supplier orders
CREATE TABLE supplier_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    supplier_name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending','shipped','received','cancelled')),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- Supplier order items
CREATE TABLE supplier_order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL,
    item_id UUID NOT NULL,
    quantity INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES supplier_orders(id),
    FOREIGN KEY (item_id) REFERENCES items(id)
);

-- Create indexes for better performance
CREATE INDEX idx_items_tenant_id ON items(tenant_id);
CREATE INDEX idx_items_sku ON items(sku);
CREATE INDEX idx_stock_movements_tenant_id ON stock_movements(tenant_id);
CREATE INDEX idx_stock_movements_item_id ON stock_movements(item_id);
CREATE INDEX idx_supplier_orders_tenant_id ON supplier_orders(tenant_id);
CREATE INDEX idx_supplier_order_items_order_id ON supplier_order_items(order_id);
CREATE INDEX idx_supplier_order_items_item_id ON supplier_order_items(item_id);
