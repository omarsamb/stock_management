package handlers

import (
	"net/http"
	"stock_management/dto"
	"stock_management/models"
	"stock_management/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ArticleHandler struct {
	Service *services.ArticleService
}

func NewArticleHandler(s *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{Service: s}
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)
	role := c.GetString("role")
	userShopIDStr := c.GetString("shop_id")

	// Determine ShopID for initial stock
	var shopID *uuid.UUID
	// If vendor, force their shop
	if role == "vendor" && userShopIDStr != "" {
		id, err := uuid.Parse(userShopIDStr)
		if err == nil {
			shopID = &id
		}
	} else if req.ShopID != nil {
		// If admin/owner and they selected a shop
		shopID = req.ShopID
	}

	article := &models.Article{
		AccountID:    accountID,
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		CategoryID:   req.CategoryID,
		BrandID:      req.BrandID,
		MinThreshold: req.MinThreshold,
		Price:        req.Price,
	}

	if err := h.Service.CreateArticle(article, req.InitialStock, shopID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	var shopID *uuid.UUID

	userRole := c.GetString("role")
	userShopIDStr := c.GetString("shop_id")

	// If user is specialized vendor (has shop_id and is vendor/manager), force shop_id
	// Actually, generic way: if user is logged in as vendor, enforce their shop
	// If user is specialized vendor (has shop_id and is vendor/manager), force shop_id
	// Actually, generic way: if user is logged in as vendor, enforce their shop
	if userRole == "vendor" {
		if userShopIDStr == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "vendor account configuration error: no shop assigned or outdated token"})
			return
		}
		id, err := uuid.Parse(userShopIDStr)
		if err == nil {
			shopID = &id
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid shop id in context"})
			return
		}
	} else {
		// Otherwise (Admin or no specific restriction yet), check query param
		shopIDStr := c.Query("shop_id")
		if shopIDStr != "" {
			id, err := uuid.Parse(shopIDStr)
			if err == nil {
				shopID = &id
			}
		}
	}

	articles, err := h.Service.GetArticlesByAccount(accountID, shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) ImportArticles(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	count, err := h.Service.ImportArticlesFromCSV(accountID, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "processed": count})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Import successful",
		"imported": count,
	})
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	articleIDStr := c.Param("id")
	articleID, _ := uuid.Parse(articleIDStr)

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	var article models.Article
	if err := h.Service.DB.Where("id = ? AND account_id = ?", articleID, accountID).First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	article.Name = req.Name
	article.Description = req.Description
	article.MinThreshold = req.MinThreshold
	article.Price = req.Price

	if err := h.Service.UpdateArticle(&article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}
