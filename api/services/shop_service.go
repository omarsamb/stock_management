package services

import (
	"stock_management/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopService struct {
	DB *gorm.DB
}

func NewShopService(db *gorm.DB) *ShopService {
	return &ShopService{DB: db}
}

func (s *ShopService) CreateShop(accountID uuid.UUID, name, location string) (*models.Shop, error) {
	shop := &models.Shop{
		ID:        uuid.New(),
		AccountID: accountID,
		Name:      name,
		Location:  location,
	}
	err := s.DB.Create(shop).Error
	return shop, err
}

func (s *ShopService) GetShopsByAccount(accountID uuid.UUID) ([]models.Shop, error) {
	var shops []models.Shop
	err := s.DB.Where("account_id = ?", accountID).Find(&shops).Error
	return shops, err
}
