package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountStatus string

const (
	AccountStatusActive    AccountStatus = "active"
	AccountStatusTrial     AccountStatus = "trial"
	AccountStatusReadOnly  AccountStatus = "read_only"
	AccountStatusSuspended AccountStatus = "suspended"
)

type Account struct {
	ID                    uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CompanyName           string         `gorm:"not null" json:"company_name"`
	SubscriptionPlan      string         `gorm:"not null" json:"subscription_plan"` // basic, pro, premium
	Status                AccountStatus  `gorm:"default:'trial'" json:"status"`
	PrimaryColor          string         `gorm:"default:'#4f46e5'" json:"primary_color"`
	BackgroundImage       string         `json:"background_image"`
	SubscriptionExpiresAt *time.Time     `json:"subscription_expires_at"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
