package models

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Tenant represents a tenant (shop) in the system
type Tenant struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never include in JSON responses
	TenantID     uuid.UUID `json:"tenant_id" db:"tenant_id"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// Item represents an inventory item
type Item struct {
	ID          uuid.UUID `json:"id" db:"id"`
	TenantID    uuid.UUID `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name"`
	SKU         *string   `json:"sku" db:"sku"` // Nullable
	Quantity    int       `json:"quantity" db:"quantity"`
	MinQuantity int       `json:"min_quantity" db:"min_quantity"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// StockMovement represents a change in stock quantity
type StockMovement struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TenantID  uuid.UUID `json:"tenant_id" db:"tenant_id"`
	ItemID    uuid.UUID `json:"item_id" db:"item_id"`
	Change    int       `json:"change" db:"change"` // Positive or negative
	Reason    *string   `json:"reason" db:"reason"` // Nullable
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// SupplierOrder represents an order from a supplier
type SupplierOrder struct {
	ID           uuid.UUID `json:"id" db:"id"`
	TenantID     uuid.UUID `json:"tenant_id" db:"tenant_id"`
	SupplierName string    `json:"supplier_name" db:"supplier_name"`
	Status       string    `json:"status" db:"status"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// SupplierOrderItem represents an item within a supplier order
type SupplierOrderItem struct {
	ID       uuid.UUID `json:"id" db:"id"`
	OrderID  uuid.UUID `json:"order_id" db:"order_id"`
	ItemID   uuid.UUID `json:"item_id" db:"item_id"`
	Quantity int       `json:"quantity" db:"quantity"`
}

// UserRole constants
const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleStaff   = "staff"
)

// SupplierOrderStatus constants
const (
	StatusPending   = "pending"
	StatusShipped   = "shipped"
	StatusReceived  = "received"
	StatusCancelled = "cancelled"
)

// SignupRequest represents the request body for user signup
type SignupRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	TenantName string `json:"tenant_name"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the response for authentication endpoints
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Role     string    `json:"role"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}