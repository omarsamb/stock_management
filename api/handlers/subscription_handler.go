package handlers

import (
	"net/http"
	"stock_management/dto"
	"stock_management/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	Service        *services.SubscriptionService
	AccountService *services.AccountService
}

func NewSubscriptionHandler(s *services.SubscriptionService, a *services.AccountService) *SubscriptionHandler {
	return &SubscriptionHandler{Service: s, AccountService: a}
}

// SelectPlan updates the desired plan and potentially returns a payment link (mock)
func (h *SubscriptionHandler) SelectPlan(c *gin.Context) {
	var req dto.SelectPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	if err := h.AccountService.UpdatePlan(accountID, req.Plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// In a real app, you would redirect to PayDunya or return a payment URL
	c.JSON(http.StatusOK, gin.H{
		"message":     "Plan selected. Please proceed to payment.",
		"payment_url": "https://paydunya.com/checkout/mock-token",
	})
}

// HandlePayDunyaWebhook handles incoming payment notifications from PayDunya
func (h *SubscriptionHandler) HandlePayDunyaWebhook(c *gin.Context) {
	// PayDunya usually sends data in POST
	var payload struct {
		Data struct {
			Status        string  `json:"status"`
			Description   string  `json:"description"`
			Amount        float64 `json:"amount"`
			CustomData    string  `json:"custom_data"` // We'll put account_id here
			TransactionID string  `json:"token"`
		} `json:"data"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Verify status
	if payload.Data.Status != "completed" {
		c.JSON(http.StatusOK, gin.H{"status": "ignored"})
		return
	}

	accountID, err := uuid.Parse(payload.Data.CustomData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid custom data (account_id)"})
		return
	}

	err = h.Service.ConfirmPayment(accountID, payload.Data.TransactionID, payload.Data.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
