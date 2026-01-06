package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Supplier struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID uuid.UUID      `gorm:"type:uuid;not null;index" json:"account_id"`
	Name      string         `gorm:"not null" json:"name"`
	Contact   string         `json:"contact"`
	Email     string         `json:"email"`
	Address   string         `json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Account Account `gorm:"foreignKey:AccountID" json:"-"`
}

type OrderStatus string

const (
	OrderDraft     OrderStatus = "draft"
	OrderSent      OrderStatus = "sent"
	OrderPartial   OrderStatus = "partial"
	OrderReceived  OrderStatus = "received"
	OrderCancelled OrderStatus = "cancelled"
)

type PurchaseOrder struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID   uuid.UUID   `gorm:"type:uuid;not null;index" json:"account_id"`
	SupplierID  uuid.UUID   `gorm:"type:uuid;not null;index" json:"supplier_id"`
	Status      OrderStatus `gorm:"default:'draft'" json:"status"`
	TotalAmount float64     `json:"total_amount"`
	OrderDate   time.Time   `json:"order_date"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`

	Account  Account  `gorm:"foreignKey:AccountID" json:"-"`
	Supplier Supplier `gorm:"foreignKey:SupplierID" json:"-"`
}

type PurchaseOrderItem struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	PurchaseOrderID uuid.UUID `gorm:"type:uuid;not null;index" json:"purchase_order_id"`
	ArticleID       uuid.UUID `gorm:"type:uuid;not null;index" json:"article_id"`
	Quantity        int       `gorm:"not null" json:"quantity"`
	ReceivedQty     int       `gorm:"default:0" json:"received_qty"`
	UnitPrice       float64   `json:"unit_price"`

	Article Article `gorm:"foreignKey:ArticleID" json:"-"`
}
