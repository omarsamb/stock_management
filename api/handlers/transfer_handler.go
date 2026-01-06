package handlers

import (
	"net/http"
	"stock_management/dto"
	"stock_management/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransferHandler struct {
	Service *services.StockService
}

func NewTransferHandler(s *services.StockService) *TransferHandler {
	return &TransferHandler{Service: s}
}

func (h *TransferHandler) InitiateTransfer(c *gin.Context) {
	var req dto.TransferStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	deviceID := req.DeviceID
	if deviceID == "" {
		deviceID = c.GetHeader("X-Device-ID")
		if deviceID == "" {
			deviceID = c.Request.UserAgent()
		}
	}

	transfer, err := h.Service.InitiateTransfer(
		accountID, req.FromShopID, req.ToShopID, req.ArticleID, userID,
		req.Qty, req.Reason, deviceID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transfer)
}

func (h *TransferHandler) ReceiveTransfer(c *gin.Context) {
	transferIDStr := c.Param("id")
	transferID, err := uuid.Parse(transferIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transfer id"})
		return
	}

	var req struct {
		DeviceID string `json:"device_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Ignore error if body is empty
	}

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	deviceID := req.DeviceID
	if deviceID == "" {
		deviceID = c.GetHeader("X-Device-ID")
		if deviceID == "" {
			deviceID = c.Request.UserAgent()
		}
	}

	if err := h.Service.ReceiveTransfer(accountID, transferID, userID, deviceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transfer received successfully"})
}

func (h *TransferHandler) ListTransfers(c *gin.Context) {
	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	transfers, err := h.Service.GetTransfers(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transfers)
}
