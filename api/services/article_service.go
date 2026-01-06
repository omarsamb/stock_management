package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"stock_management/models"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArticleService struct {
	DB *gorm.DB
}

func NewArticleService(db *gorm.DB) *ArticleService {
	return &ArticleService{DB: db}
}

func (s *ArticleService) CreateArticle(article *models.Article) error {
	if article.SKU == "" {
		sku, err := s.GenerateSKU(article.AccountID, "")
		if err != nil {
			return err
		}
		article.SKU = sku
	}
	return s.DB.Create(article).Error
}

func (s *ArticleService) GenerateSKU(accountID uuid.UUID, category string) (string, error) {
	// Generate a unique SKU: CAT-XXXXX or ART-XXXXX
	prefix := "ART"
	if category != "" {
		prefix = strings.ToUpper(category[:3])
	}

	for i := 0; i < 5; i++ { // Try 5 times to find a unique SKU
		randomPart := fmt.Sprintf("%05d", rand.Intn(100000))
		sku := fmt.Sprintf("%s-%s", prefix, randomPart)

		var count int64
		s.DB.Model(&models.Article{}).Where("account_id = ? AND sku = ?", accountID, sku).Count(&count)
		if count == 0 {
			return sku, nil
		}
	}

	return "", fmt.Errorf("could not generate a unique SKU after several attempts")
}

func (s *ArticleService) GetArticlesByAccount(accountID uuid.UUID) ([]models.Article, error) {
	var articles []models.Article
	// We use a subquery to sum up the stock levels across all shops for each article
	err := s.DB.Model(&models.Article{}).
		Select("articles.*, (SELECT COALESCE(SUM(quantity), 0) FROM stock_levels WHERE stock_levels.article_id = articles.id) as total_stock").
		Where("account_id = ?", accountID).
		Find(&articles).Error
	return articles, err
}

func (s *ArticleService) UpdateArticle(article *models.Article) error {
	return s.DB.Save(article).Error
}

func (s *ArticleService) ImportArticlesFromCSV(accountID uuid.UUID, reader io.Reader) (int, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	// Read header
	headers, err := csvReader.Read()
	if err != nil {
		return 0, err
	}

	// Simple header mapping (SKU, Name, Description, MinThreshold)
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[strings.ToLower(h)] = i
	}

	count := 0
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return count, err
		}

		article := &models.Article{
			AccountID: accountID,
			Name:      record[headerMap["name"]],
		}

		if idx, ok := headerMap["sku"]; ok {
			article.SKU = record[idx]
		}
		if idx, ok := headerMap["description"]; ok {
			article.Description = record[idx]
		}

		if err := s.CreateArticle(article); err != nil {
			// In a real app, we might want to collect errors and continue
			return count, err
		}
		count++
	}

	return count, nil
}
