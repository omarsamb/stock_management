package middleware

import (
	"net/http"
	"stock_management/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EnforceSubscriptionMiddleware blocks non-GET requests if the account is in read-only mode.
func EnforceSubscriptionMiddleware(gormDB *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID := c.GetString("account_id")
		if accountID == "" {
			c.Next()
			return
		}

		// Skip check for GET requests
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		var account models.Account
		if err := gormDB.First(&account, "id = ?", accountID).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account not found"})
			c.Abort()
			return
		}

		if account.Status == models.AccountStatusReadOnly || account.Status == models.AccountStatusSuspended {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Account is in read-only mode. Please renew your subscription.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
