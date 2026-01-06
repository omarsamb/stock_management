package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID     uuid.UUID `gorm:"type:uuid;not null;index" json:"account_id"`
	ExternalPayID string    `json:"external_pay_id"` // Reference to Mobile Money transaction
	Amount        float64   `json:"amount"`
	ExpiryDate    time.Time `json:"expiry_date"`
	CreatedAt     time.Time `json:"created_at"`

	Account Account `gorm:"foreignKey:AccountID" json:"-"`
}
