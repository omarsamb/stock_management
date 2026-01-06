package services

import (
	"gorm.io/gorm"
)

type ServicesManager struct {
	DB                  *gorm.DB
	JWTSecret           string
	StockService        *StockService
	ArticleService      *ArticleService
	AccountService      *AccountService
	SubscriptionService *SubscriptionService
	ShopService         *ShopService
	WhatsAppService     *WhatsAppService
}

func InitServices(db *gorm.DB, jwtSecret string) *ServicesManager {
	return &ServicesManager{
		DB:                  db,
		JWTSecret:           jwtSecret,
		StockService:        NewStockService(db),
		ArticleService:      NewArticleService(db),
		AccountService:      NewAccountService(db, jwtSecret),
		SubscriptionService: NewSubscriptionService(db),
		ShopService:         NewShopService(db),
		WhatsAppService:     NewWhatsAppService(),
	}
}
