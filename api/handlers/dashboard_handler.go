package handlers

import (
	"net/http"
	"stock_management/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DashboardHandler struct {
	Service *services.StockService
}

func NewDashboardHandler(s *services.StockService) *DashboardHandler {
	return &DashboardHandler{Service: s}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	shopID := uuid.Nil
	if shopIDStr := c.Query("shop_id"); shopIDStr != "" {
		shopID, _ = uuid.Parse(shopIDStr)
	}

	stats, err := h.Service.GetDashboardStats(accountID, shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
