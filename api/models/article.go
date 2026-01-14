package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Article struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"account_id"`
	Code         string         `gorm:"not null;index" json:"code"`
	Name         string         `gorm:"not null" json:"name"`
	Description  string         `json:"description"`
	CategoryID   *uuid.UUID     `gorm:"type:uuid;index" json:"category_id"`
	BrandID      *uuid.UUID     `gorm:"type:uuid;index" json:"brand_id"`
	MinThreshold int            `gorm:"default:0" json:"min_threshold"`
	Price        float64        `gorm:"type:decimal(10,2);default:0" json:"price"`
	TotalStock   int            `gorm:"->" json:"total_stock"`
	ImageURL     string         `json:"image_url"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Account Account `gorm:"foreignKey:AccountID" json:"-"`
}

func (a *Article) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"account_id"`
	Name      string    `gorm:"not null" json:"name"`
}

type Brand struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"account_id"`
	Name      string    `gorm:"not null" json:"name"`
}
