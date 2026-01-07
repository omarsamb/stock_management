package services

import (
	"errors"
	"stock_management/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockService struct {
	DB *gorm.DB
}

func NewStockService(db *gorm.DB) *StockService {
	return &StockService{DB: db}
}

// RecordMovement registers a stock movement and updates the stock level in a transaction.
func (s *StockService) RecordMovement(
	accountID, shopID, articleID, userID uuid.UUID,
	moveType models.MovementType,
	qty int,
	reason, deviceID string,
) (*models.StockMovement, error) {
	var movement *models.StockMovement

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Get or create current stock level
		var stock models.StockLevel
		res := tx.Where("article_id = ? AND shop_id = ?", articleID, shopID).First(&stock)

		oldQty := 0
		if res.Error == nil {
			oldQty = stock.Quantity
		} else if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			stock = models.StockLevel{
				ArticleID: articleID,
				ShopID:    shopID,
				Quantity:  0,
			}
		} else {
			return res.Error
		}

		// 2. Calculate new quantity
		newQty := oldQty
		switch moveType {
		case models.MovementIn:
			newQty += qty
		case models.MovementOut:
			if oldQty < qty {
				return errors.New("insufficient stock")
			}
			newQty -= qty
		case models.MovementAdjust:
			newQty = qty
		case models.MovementTransfer:
			// Transfer is handled separately as it involves two shops
			return errors.New("use TransferStock for transfers")
		}

		// 3. Update stock level
		stock.Quantity = newQty
		if err := tx.Save(&stock).Error; err != nil {
			return err
		}

		// 4. Create movement log (Audit Log)
		movement = &models.StockMovement{
			ID:        uuid.New(),
			AccountID: accountID,
			ShopID:    shopID,
			ArticleID: articleID,
			UserID:    userID,
			Type:      moveType,
			Qty:       qty,
			OldValue:  oldQty,
			NewValue:  newQty,
			Reason:    reason,
			DeviceID:  deviceID,
		}
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		return nil
	})

	return movement, err
}

func (s *StockService) InitiateTransfer(
	accountID, fromShopID, toShopID, articleID, userID uuid.UUID,
	qty int,
	reason, deviceID string,
) (*models.StockTransfer, error) {
	var transfer *models.StockTransfer

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Exit from source shop (immediate)
		service := NewStockService(tx)
		_, err := service.RecordMovement(accountID, fromShopID, articleID, userID, models.MovementOut, qty, "Transfer Out: "+reason, deviceID)
		if err != nil {
			return err
		}

		// 2. Create Transfer Record (Pending)
		transfer = &models.StockTransfer{
			AccountID:   accountID,
			FromShopID:  fromShopID,
			ToShopID:    toShopID,
			ArticleID:   articleID,
			Qty:         qty,
			Status:      models.TransferStatusPending,
			InitiatedBy: userID,
		}
		if err := tx.Create(transfer).Error; err != nil {
			return err
		}

		return nil
	})

	return transfer, err
}

func (s *StockService) ReceiveTransfer(
	accountID, transferID, userID uuid.UUID,
	deviceID string,
) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var transfer models.StockTransfer
		if err := tx.First(&transfer, "id = ? AND account_id = ?", transferID, accountID).Error; err != nil {
			return err
		}

		if transfer.Status != models.TransferStatusPending {
			return errors.New("transfer is not in pending status")
		}

		// 1. Entry to destination shop
		service := NewStockService(tx)
		_, err := service.RecordMovement(accountID, transfer.ToShopID, transfer.ArticleID, userID, models.MovementIn, transfer.Qty, "Transfer In (Received)", deviceID)
		if err != nil {
			return err
		}

		// 2. Update Transfer Record
		now := time.Now()
		transfer.Status = models.TransferStatusReceived
		transfer.ReceivedBy = &userID
		transfer.ReceivedAt = &now

		if err := tx.Save(&transfer).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *StockService) GetTransfers(accountID uuid.UUID) ([]models.StockTransfer, error) {
	var transfers []models.StockTransfer
	err := s.DB.Preload("FromShop").Preload("ToShop").Preload("Article").
		Where("account_id = ?", accountID).
		Order("created_at desc").Find(&transfers).Error
	return transfers, err
}

func (s *StockService) GetStockLevels(accountID, shopID uuid.UUID) ([]models.StockLevel, error) {
	var levels []models.StockLevel
	err := s.DB.Preload("Article").Where("account_id = ? AND shop_id = ?", accountID, shopID).Find(&levels).Error
	return levels, err
}

func (s *StockService) GetMovements(accountID, shopID, articleID uuid.UUID) ([]models.StockMovement, error) {
	var movements []models.StockMovement
	query := s.DB.Where("account_id = ?", accountID)
	if shopID != uuid.Nil {
		query = query.Where("shop_id = ?", shopID)
	}
	if articleID != uuid.Nil {
		query = query.Where("article_id = ?", articleID)
	}
	err := query.Order("created_at desc").Limit(20).Find(&movements).Error
	return movements, err
}

type DashboardStats struct {
	TotalStockValue float64 `json:"total_stock_value"`
	LowStockAlerts  int64   `json:"low_stock_alerts"`
	TotalArticles   int64   `json:"total_articles"`
	ActiveShops     int64   `json:"active_shops"`
}

func (s *StockService) GetDashboardStats(accountID, shopID uuid.UUID) (*DashboardStats, error) {
	var stats DashboardStats

	// 1. Total Articles
	s.DB.Model(&models.Article{}).Where("account_id = ?", accountID).Count(&stats.TotalArticles)

	// 2. Active Shops
	s.DB.Model(&models.Shop{}).Where("account_id = ?", accountID).Count(&stats.ActiveShops)

	// 3. Low Stock Alerts
	query := s.DB.Table("stock_levels").
		Joins("JOIN articles ON articles.id = stock_levels.article_id").
		Where("articles.account_id = ? AND stock_levels.quantity < articles.min_threshold", accountID)
	if shopID != uuid.Nil {
		query = query.Where("stock_levels.shop_id = ?", shopID)
	}
	query.Count(&stats.LowStockAlerts)

	// 4. Total Stock Value
	var totalValue float64
	queryValue := s.DB.Table("stock_levels").
		Joins("JOIN articles ON articles.id = stock_levels.article_id").
		Where("articles.account_id = ?", accountID)
	if shopID != uuid.Nil {
		queryValue = queryValue.Where("stock_levels.shop_id = ?", shopID)
	}
	queryValue.Select("SUM(stock_levels.quantity * articles.price)").Row().Scan(&totalValue)

	stats.TotalStockValue = totalValue

	return &stats, nil
}
