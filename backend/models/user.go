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

// UserRole constants
const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleStaff   = "staff"
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