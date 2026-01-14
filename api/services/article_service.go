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

func (s *ArticleService) CreateArticle(article *models.Article, initialStock int, shopID *uuid.UUID, userID uuid.UUID) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if article.Code == "" {
			code, err := s.GenerateCode(article.AccountID, "")
			if err != nil {
				return err
			}
			article.Code = code
		}
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		// Add Initial Stock if provided and shop is specified
		if initialStock > 0 && shopID != nil && *shopID != uuid.Nil {
			stockService := NewStockService(tx)
			_, err := stockService.RecordMovement(article.AccountID, *shopID, article.ID, userID, models.MovementIn, initialStock, "Initial Stock", "system")
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *ArticleService) GenerateCode(accountID uuid.UUID, category string) (string, error) {
	// Generate a unique Code: CAT-XXXXX or ART-XXXXX
	prefix := "ART"
	if category != "" {
		prefix = strings.ToUpper(category[:3])
	}

	for i := 0; i < 5; i++ { // Try 5 times to find a unique Code
		randomPart := fmt.Sprintf("%05d", rand.Intn(100000))
		code := fmt.Sprintf("%s-%s", prefix, randomPart)

		var count int64
		s.DB.Model(&models.Article{}).Where("account_id = ? AND code = ?", accountID, code).Count(&count)
		if count == 0 {
			return code, nil
		}
	}

	return "", fmt.Errorf("could not generate a unique Code after several attempts")
}

func (s *ArticleService) GetArticlesByAccount(accountID uuid.UUID, shopID *uuid.UUID) ([]models.Article, error) {
	var articles []models.Article

	selectQuery := "articles.*"
	stockSubQuery := "(SELECT COALESCE(SUM(quantity), 0) FROM stock_levels WHERE stock_levels.article_id = articles.id"

	if shopID != nil {
		stockSubQuery += fmt.Sprintf(" AND stock_levels.shop_id = '%s'", shopID.String())
	}

	stockSubQuery += ") as total_stock"

	query := s.DB.Model(&models.Article{}).
		Select(selectQuery+", "+stockSubQuery).
		Where("articles.account_id = ?", accountID)

	if shopID != nil {
		query = query.Joins("JOIN stock_levels ON stock_levels.article_id = articles.id").
			Where("stock_levels.shop_id = ?", shopID)
	}

	err := query.Find(&articles).Error
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

	// Simple header mapping (Code, Name, Description, MinThreshold)
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

		if idx, ok := headerMap["code"]; ok {
			article.Code = record[idx]
		} else if idx, ok := headerMap["sku"]; ok {
			article.Code = record[idx]
		}
		if idx, ok := headerMap["description"]; ok {
			article.Description = record[idx]
		}

		// For import, we don't set initial stock or shop for now
		if err := s.CreateArticle(article, 0, nil, uuid.Nil); err != nil {
			// In a real app, we might want to collect errors and continue
			return count, err
		}
		count++
	}

	return count, nil
}
