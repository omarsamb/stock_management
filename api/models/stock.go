package models

import (
	"time"

	"github.com/google/uuid"
)

type StockLevel struct {
	ArticleID uuid.UUID `gorm:"type:uuid;primaryKey" json:"article_id"`
	ShopID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"shop_id"`
	Quantity  int       `gorm:"default:0" json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`

	Article Article `gorm:"foreignKey:ArticleID"`
	Shop    Shop    `gorm:"foreignKey:ShopID"`
}

type MovementType string

const (
	MovementIn       MovementType = "in"       // Reception, Return
	MovementOut      MovementType = "out"      // Sale, Loss, Expired
	MovementTransfer MovementType = "transfer" // Inter-shop
	MovementAdjust   MovementType = "adjust"   // Physical count adjustment
)

type StockMovement struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID uuid.UUID    `gorm:"type:uuid;not null;index" json:"account_id"`
	ShopID    uuid.UUID    `gorm:"type:uuid;not null;index" json:"shop_id"`
	ArticleID uuid.UUID    `gorm:"type:uuid;not null;index" json:"article_id"`
	UserID    uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      MovementType `gorm:"not null" json:"type"`
	Qty       int          `gorm:"not null" json:"qty"`
	OldValue  int          `json:"old_value"`
	NewValue  int          `json:"new_value"`
	Reason    string       `json:"reason"`
	DeviceID  string       `json:"device_id"`
	CreatedAt time.Time    `json:"created_at"`

	Account Account `gorm:"foreignKey:AccountID" json:"-"`
	Shop    Shop    `gorm:"foreignKey:ShopID" json:"-"`
	Article Article `gorm:"foreignKey:ArticleID" json:"-"`
	User    User    `gorm:"foreignKey:UserID" json:"-"`
}
