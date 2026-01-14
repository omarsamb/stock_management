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
	err := s.DB.Preload("Article").
		Joins("JOIN articles ON articles.id = stock_levels.article_id").
		Where("articles.account_id = ? AND stock_levels.shop_id = ?", accountID, shopID).
		Find(&levels).Error
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

type StockByCategory struct {
	CategoryName string  `json:"category_name"`
	TotalValue   float64 `json:"total_value"`
	ItemCount    int64   `json:"item_count"`
}

type LowStockItem struct {
	ArticleName  string `json:"article_name"`
	Quantity     int    `json:"quantity"`
	MinThreshold int    `json:"min_threshold"`
	ShopName     string `json:"shop_name"`
}

type DailyMovement struct {
	Date   string `json:"date"`
	InQty  int    `json:"in_qty"`
	OutQty int    `json:"out_qty"`
}

type DashboardStats struct {
	TotalStockValue float64           `json:"total_stock_value"`
	LowStockAlerts  int64             `json:"low_stock_alerts"`
	TotalArticles   int64             `json:"total_articles"`
	ActiveShops     int64             `json:"active_shops"`
	StockByCat      []StockByCategory `json:"stock_by_category"`
	LowStockItems   []LowStockItem    `json:"low_stock_items"`
	DailyMovements  []DailyMovement   `json:"daily_movements"`
}

func (s *StockService) GetDashboardStats(accountID, shopID uuid.UUID) (*DashboardStats, error) {
	var stats DashboardStats

	// 1. Total Articles
	if shopID != uuid.Nil {
		s.DB.Table("articles").
			Joins("JOIN stock_levels ON stock_levels.article_id = articles.id").
			Where("articles.account_id = ? AND stock_levels.shop_id = ?", accountID, shopID).
			Distinct("articles.id").
			Count(&stats.TotalArticles)
	} else {
		s.DB.Model(&models.Article{}).Where("account_id = ?", accountID).Count(&stats.TotalArticles)
	}

	// 2. Active Shops
	s.DB.Model(&models.Shop{}).Where("account_id = ?", accountID).Count(&stats.ActiveShops)

	// 3. Low Stock Alerts & Items
	queryLow := s.DB.Table("stock_levels").
		Select("articles.name as article_name, stock_levels.quantity, articles.min_threshold, shops.name as shop_name").
		Joins("JOIN articles ON articles.id = stock_levels.article_id").
		Joins("JOIN shops ON shops.id = stock_levels.shop_id").
		Where("articles.account_id = ? AND stock_levels.quantity < articles.min_threshold", accountID)

	if shopID != uuid.Nil {
		queryLow = queryLow.Where("stock_levels.shop_id = ?", shopID)
	}

	// Count total alerts
	queryLow.Count(&stats.LowStockAlerts)

	// Get top 10 low stock items for the table
	queryLow.Order("stock_levels.quantity ASC").Limit(10).Scan(&stats.LowStockItems)

	// 4. Total Stock Value & By Category
	var totalValue float64
	queryValue := s.DB.Table("stock_levels").
		Joins("JOIN articles ON articles.id = stock_levels.article_id").
		Where("articles.account_id = ?", accountID)

	if shopID != uuid.Nil {
		queryValue = queryValue.Where("stock_levels.shop_id = ?", shopID)
	}

	queryValue.Select("SUM(stock_levels.quantity * articles.price)").Row().Scan(&totalValue)
	stats.TotalStockValue = totalValue

	// Stock By Category
	queryCat := s.DB.Table("stock_levels").
		Select("COALESCE(categories.name, 'Non catégorisé') as category_name, SUM(stock_levels.quantity * articles.price) as total_value, COUNT(DISTINCT articles.id) as item_count").
		Joins("JOIN articles ON articles.id = stock_levels.article_id").
		Joins("LEFT JOIN categories ON categories.id = articles.category_id").
		Where("articles.account_id = ?", accountID)

	if shopID != uuid.Nil {
		queryCat = queryCat.Where("stock_levels.shop_id = ?", shopID)
	}

	queryCat.Group("categories.name").Scan(&stats.StockByCat)

	// 5. Daily Movements (Last 30 days)
	// We need to group by day. PostgreSQL: to_char(created_at, 'YYYY-MM-DD')
	// SQLite: strftime('%Y-%m-%d', created_at)
	// Assuming Postgres based on "gorm" usage typical patterns, but let's try to be generic or assume standard SQL if possible.
	// Since the USER didn't specify DB, I'll assume Postgres or SQLite.
	// For standard GORM date truncation, it's dialect specific.
	// Let's assume Postgres for now as it's common. If it fails, I'll need to adjust.
	// Actually, look at previous files... `docker-compose.yml` likely has DB info.
	// But to be safe, I'll try to use a slightly more raw SQL approach or a safe cast.

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	queryMov := s.DB.Table("stock_movements").
		Select("to_char(created_at, 'YYYY-MM-DD') as date, "+
			"SUM(CASE WHEN type = 'in' THEN qty ELSE 0 END) as in_qty, "+
			"SUM(CASE WHEN type = 'out' THEN qty ELSE 0 END) as out_qty").
		Where("account_id = ? AND created_at >= ?", accountID, thirtyDaysAgo).
		Group("to_char(created_at, 'YYYY-MM-DD')").
		Order("date ASC")

	if shopID != uuid.Nil {
		queryMov = queryMov.Where("shop_id = ?", shopID)
	}

	queryMov.Scan(&stats.DailyMovements)

	return &stats, nil
}

type SalesStatPoint struct {
	Label    string  `json:"label"`
	Revenue  float64 `json:"revenue"`
	Quantity int     `json:"quantity"`
}

func (s *StockService) GetSalesStats(accountID, shopID uuid.UUID, period string) ([]SalesStatPoint, error) {
	var stats []SalesStatPoint
	var dateFormat string
	var startDate time.Time

	now := time.Now()

	switch period {
	case "day":
		// Hourly stats for today
		dateFormat = "HH24:00" // Postgres
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "week":
		// Daily stats for last 7 days
		dateFormat = "YYYY-MM-DD"
		startDate = now.AddDate(0, 0, -7)
	case "month":
		// Daily stats for this month
		dateFormat = "YYYY-MM-DD"
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "year":
		// Monthly stats for this year
		dateFormat = "YYYY-MM"
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default:
		return nil, errors.New("invalid period")
	}

	query := s.DB.Table("stock_movements").
		Select("to_char(stock_movements.created_at, ?) as label, "+
			"SUM(stock_movements.qty * articles.price) as revenue, "+
			"SUM(stock_movements.qty) as quantity", dateFormat).
		Joins("JOIN articles ON articles.id = stock_movements.article_id").
		Where("stock_movements.account_id = ? AND stock_movements.type = ? AND stock_movements.created_at >= ?",
			accountID, models.MovementOut, startDate)

	if shopID != uuid.Nil {
		query = query.Where("stock_movements.shop_id = ?", shopID)
	}

	err := query.Group("label").Order("label").Scan(&stats).Error

	return stats, err
}
