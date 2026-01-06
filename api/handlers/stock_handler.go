package handlers

import (
	"net/http"
	"stock_management/dto"
	"stock_management/models"
	"stock_management/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StockHandler struct {
	Service *services.StockService
}

func NewStockHandler(s *services.StockService) *StockHandler {
	return &StockHandler{Service: s}
}

func (h *StockHandler) RecordMovement(c *gin.Context) {
	var req dto.RecordMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	moveType := models.MovementType(req.Type)

	deviceID := req.DeviceID
	if deviceID == "" {
		deviceID = c.GetHeader("X-Device-ID")
		if deviceID == "" {
			deviceID = c.Request.UserAgent()
		}
	}

	movement, err := h.Service.RecordMovement(
		accountID, req.ShopID, req.ArticleID, userID,
		moveType, req.Qty, req.Reason, deviceID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movement)
}

func (h *StockHandler) ListStockLevels(c *gin.Context) {
	shopIDStr := c.Query("shop_id")
	if shopIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "shop_id is required"})
		return
	}
	shopID, _ := uuid.Parse(shopIDStr)

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	levels, err := h.Service.GetStockLevels(accountID, shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, levels)
}

func (h *StockHandler) ListMovements(c *gin.Context) {
	shopIDStr := c.Query("shop_id")
	articleIDStr := c.Query("article_id")

	var shopID, articleID uuid.UUID
	if shopIDStr != "" {
		shopID, _ = uuid.Parse(shopIDStr)
	}
	if articleIDStr != "" {
		articleID, _ = uuid.Parse(articleIDStr)
	}

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	movements, err := h.Service.GetMovements(accountID, shopID, articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movements)
}
