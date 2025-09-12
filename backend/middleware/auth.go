package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"stock_management/backend/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-change-this-in-production"
	}
	jwtSecret = []byte(secret)
}

// JWTMiddleware validates JWT tokens and adds user info to the request context
func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*models.JWTClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserFromContext extracts user claims from the request context
func GetUserFromContext(r *http.Request) (*models.JWTClaims, bool) {
	claims, ok := r.Context().Value(UserContextKey).(*models.JWTClaims)
	return claims, ok
}

// RequireRole creates a middleware that checks if the user has the required role
func RequireRole(role string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserFromContext(r)
			if !ok {
				http.Error(w, "User not found in context", http.StatusUnauthorized)
				return
			}

			if claims.Role != role && claims.Role != models.RoleAdmin {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireMinimumRole creates a middleware that checks if the user has at least the minimum required role
func RequireMinimumRole(minRole string) func(http.HandlerFunc) http.HandlerFunc {
	roleHierarchy := map[string]int{
		models.RoleStaff:   1,
		models.RoleManager: 2,
		models.RoleAdmin:   3,
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserFromContext(r)
			if !ok {
				http.Error(w, "User not found in context", http.StatusUnauthorized)
				return
			}

			userRoleLevel := roleHierarchy[claims.Role]
			minRoleLevel := roleHierarchy[minRole]

			if userRoleLevel < minRoleLevel {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GenerateJWT creates a new JWT token for the user
func GenerateJWT(user *models.User) (string, error) {
	claims := &models.JWTClaims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Role:     user.Role,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}