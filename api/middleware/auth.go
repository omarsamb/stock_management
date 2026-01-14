package middleware

import (
	"net/http"
	"strings"

	"stock_management/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles JWT authentication and tenant isolation.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Inject claims into context
		c.Set("account_id", claims.AccountID)
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		if claims.ShopID != "" {
			c.Set("shop_id", claims.ShopID)
		}

		c.Next()
	}
}
