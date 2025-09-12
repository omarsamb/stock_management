# Stock Management API Documentation

This document describes the authentication and authorization endpoints for the Stock Management system.

## Authentication Endpoints

### POST /api/signup
Creates a new tenant and admin user.

**Request Body:**
```json
{
  "email": "admin@company.com",
  "password": "password123",
  "tenant_name": "My Company"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "admin@company.com",
    "tenant_id": "uuid",
    "role": "admin",
    "created_at": "2025-09-12T08:16:11.840861Z"
  }
}
```

### POST /api/login
Authenticates a user and returns a JWT token.

**Request Body:**
```json
{
  "email": "user@company.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "user@company.com",
    "tenant_id": "uuid",
    "role": "manager",
    "created_at": "2025-09-12T08:16:11.840861Z"
  }
}
```

### POST /api/logout
Logs out the user (client-side token removal).

**Response:**
```json
{
  "message": "Logged out successfully"
}
```

## Protected Endpoints

All protected endpoints require an `Authorization: Bearer <token>` header.

### GET /api/profile
Returns the current user's profile information.

**Response:**
```json
{
  "id": "uuid",
  "email": "user@company.com",
  "tenant_id": "uuid",
  "role": "manager",
  "created_at": "2025-09-12T08:16:11.840861Z"
}
```

## User Management Endpoints

### POST /api/users (Admin Only)
Creates a new user in the current tenant.

**Request Body:**
```json
{
  "email": "newuser@company.com",
  "password": "password123",
  "role": "staff"
}
```

**Response:**
```json
{
  "id": "uuid",
  "email": "newuser@company.com",
  "tenant_id": "uuid",
  "role": "staff",
  "created_at": "2025-09-12T08:16:11.840861Z"
}
```

### GET /api/users (Manager+ Only)
Lists all users in the current tenant.

**Response:**
```json
[
  {
    "id": "uuid",
    "email": "admin@company.com",
    "tenant_id": "uuid",
    "role": "admin",
    "created_at": "2025-09-12T08:16:11.840861Z"
  },
  {
    "id": "uuid",
    "email": "manager@company.com",
    "tenant_id": "uuid",
    "role": "manager",
    "created_at": "2025-09-12T08:16:12.840861Z"
  }
]
```

## Role-Based Demo Endpoints

### GET /api/staff-area (Staff+ Access)
Accessible to staff, manager, and admin roles.

### GET /api/manager-area (Manager+ Access)
Accessible to manager and admin roles only.

### GET /api/admin-area (Admin Only)
Accessible to admin role only.

## Role Hierarchy

1. **Staff** - Basic access level
   - Can access staff areas
   - Can view own profile

2. **Manager** - Mid-level access
   - Can access staff and manager areas
   - Can view all users in tenant
   - Can manage items and orders (future)

3. **Admin** - Full access
   - Can access all areas
   - Can create and manage users
   - Can manage tenant settings

## Multi-Tenancy

- Each user belongs to exactly one tenant
- Users can only see/interact with data from their own tenant
- First user created during signup becomes admin of new tenant
- All subsequent users must be created by an admin

## JWT Token

- Tokens expire after 24 hours
- Contains user_id, tenant_id, role, and email claims
- Must be included in Authorization header as "Bearer <token>"

## Error Responses

- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Insufficient permissions
- `400 Bad Request` - Invalid request body
- `404 Not Found` - Resource not found
- `409 Conflict` - User already exists (signup)
- `500 Internal Server Error` - Server error

## Environment Variables

- `DB_HOST` - Database host (default: localhost)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: stock_management)
- `JWT_SECRET` - JWT signing secret (default: development key)