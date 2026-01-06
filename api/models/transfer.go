package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransferStatus string

const (
	TransferStatusPending   TransferStatus = "pending"
	TransferStatusReceived  TransferStatus = "received"
	TransferStatusCancelled TransferStatus = "cancelled"
)

type StockTransfer struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"account_id"`
	FromShopID uuid.UUID      `gorm:"type:uuid;not null;index" json:"from_shop_id"`
	ToShopID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"to_shop_id"`
	ArticleID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"article_id"`
	Qty        int            `gorm:"not null" json:"qty"`
	Status     TransferStatus `gorm:"not null;default:'pending'" json:"status"`

	InitiatedBy uuid.UUID  `gorm:"type:uuid;not null" json:"initiated_by"`
	ReceivedBy  *uuid.UUID `gorm:"type:uuid" json:"received_by,omitempty"`

	CreatedAt  time.Time      `json:"created_at"`
	ReceivedAt *time.Time     `json:"received_at,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Account  Account `gorm:"foreignKey:AccountID" json:"-"`
	FromShop Shop    `gorm:"foreignKey:FromShopID" json:"-"`
	ToShop   Shop    `gorm:"foreignKey:ToShopID" json:"-"`
	Article  Article `gorm:"foreignKey:ArticleID" json:"-"`
}

func (t *StockTransfer) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
