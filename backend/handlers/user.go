package handlers

import (
	"encoding/json"
	"net/http"
	"stock_management/backend/database"
	"stock_management/backend/middleware"
	"stock_management/backend/models"
)

type UserHandler struct {
	db *database.DB
}

func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{db: db}
}

// CreateUser creates a new user (admin only)
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Get the current user from context
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse request
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.Role == "" {
		http.Error(w, "Email, password, and role are required", http.StatusBadRequest)
		return
	}

	// Validate role
	if req.Role != models.RoleAdmin && req.Role != models.RoleManager && req.Role != models.RoleStaff {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	// Create user in the same tenant as the admin
	user, err := h.db.CreateUser(req.Email, req.Password, claims.TenantID, req.Role)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetUsers lists all users in the current tenant (admin and manager only)
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	users, err := h.db.GetUsersByTenant(claims.TenantID)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// StaffOnlyEndpoint demonstrates staff-level access
func (h *UserHandler) StaffOnlyEndpoint(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"message":   "Welcome to the staff area",
		"user_role": claims.Role,
		"access":    "This endpoint is accessible to all authenticated users (staff, manager, admin)",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ManagerOnlyEndpoint demonstrates manager-level access
func (h *UserHandler) ManagerOnlyEndpoint(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"message":   "Welcome to the manager area",
		"user_role": claims.Role,
		"access":    "This endpoint is accessible to managers and admins only",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AdminOnlyEndpoint demonstrates admin-level access
func (h *UserHandler) AdminOnlyEndpoint(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"message":   "Welcome to the admin area",
		"user_role": claims.Role,
		"access":    "This endpoint is accessible to admins only",
		"tenant_id": claims.TenantID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}