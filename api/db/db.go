package db

import (
	"stock_management/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDatabase initializes the database connection and migrates the schema.
func InitDatabase(dbPath string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migration automatique de tous les modèles
	err = db.AutoMigrate(
		&models.Account{},
		&models.User{},
		&models.Shop{},
		&models.Article{}, &models.Category{}, &models.Brand{},
		&models.StockLevel{}, &models.StockMovement{},
		&models.Subscription{}, &models.Supplier{},
		&models.PurchaseOrder{}, &models.PurchaseOrderItem{},
		&models.StockTransfer{},
	)
	if err != nil {
		return nil, err
	}
	// Seed idempotent des données nécessaires (rôles, etc.)
	if err := SeedInitialData(db); err != nil {
		return nil, err
	}

	return db, nil
}

func SeedInitialData(db *gorm.DB) error {
	return nil
}
