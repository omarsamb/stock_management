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

	article := &models.Article{
		AccountID:    accountID,
		Name:         req.Name,
		SKU:          req.SKU,
		Description:  req.Description,
		CategoryID:   req.CategoryID,
		BrandID:      req.BrandID,
		MinThreshold: req.MinThreshold,
		Price:        req.Price,
	}

	if err := h.Service.CreateArticle(article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	articles, err := h.Service.GetArticlesByAccount(accountID)
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
