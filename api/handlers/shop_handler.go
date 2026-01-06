package handlers

import (
	"net/http"
	"stock_management/dto"
	"stock_management/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShopHandler struct {
	Service *services.ShopService
}

func NewShopHandler(s *services.ShopService) *ShopHandler {
	return &ShopHandler{Service: s}
}

func (h *ShopHandler) CreateShop(c *gin.Context) {
	var req dto.CreateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	shop, err := h.Service.CreateShop(accountID, req.Name, req.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shop)
}

func (h *ShopHandler) ListShops(c *gin.Context) {
	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	shops, err := h.Service.GetShopsByAccount(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shops)
}
