-- Drop indexes
DROP INDEX IF EXISTS idx_supplier_order_items_item_id;
DROP INDEX IF EXISTS idx_supplier_order_items_order_id;
DROP INDEX IF EXISTS idx_supplier_orders_tenant_id;
DROP INDEX IF EXISTS idx_stock_movements_item_id;
DROP INDEX IF EXISTS idx_stock_movements_tenant_id;
DROP INDEX IF EXISTS idx_items_sku;
DROP INDEX IF EXISTS idx_items_tenant_id;

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS supplier_order_items;
DROP TABLE IF EXISTS supplier_orders;
DROP TABLE IF EXISTS stock_movements;
DROP TABLE IF EXISTS items;
