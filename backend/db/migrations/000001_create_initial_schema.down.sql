-- Drop indexes
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_tenant_id;

-- Drop tables in reverse order
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;

-- Drop extension (commented out to avoid affecting other databases)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
