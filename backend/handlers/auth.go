package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"stock_management/backend/database"
	"stock_management/backend/middleware"
	"stock_management/backend/models"
)

type AuthHandler struct {
	db *database.DB
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Signup creates a new tenant and admin user
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.TenantName == "" {
		http.Error(w, "Email, password, and tenant name are required", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	_, err := h.db.GetUserByEmail(req.Email)
	if err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	if err != sql.ErrNoRows {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Create tenant
	tenant, err := h.db.CreateTenant(req.TenantName)
	if err != nil {
		http.Error(w, "Failed to create tenant", http.StatusInternalServerError)
		return
	}

	// Create admin user
	user, err := h.db.CreateUser(req.Email, req.Password, tenant.ID, models.RoleAdmin)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return response
	response := models.AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Get user by email
	user, err := h.db.GetUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Validate password
	if !h.db.ValidatePassword(user, req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return response
	response := models.AuthResponse{
		Token: token,
		User:  *user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Logout handles user logout (in a stateless JWT system, this is mainly for client-side cleanup)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// In a JWT-based system, logout is typically handled client-side by removing the token
	// For additional security, you could implement a token blacklist here
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// Profile returns the current user's profile information
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Get fresh user data from database
	user, err := h.db.GetUserByEmail(claims.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}