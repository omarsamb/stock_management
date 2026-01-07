package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleOwner   UserRole = "owner"
	RoleManager UserRole = "manager"
	RoleAnalyst UserRole = "analyst"
	RoleAdmin   UserRole = "admin"
	RoleVendor  UserRole = "vendor"
)

type User struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"account_id"`
	Phone               string         `gorm:"uniqueIndex;not null" json:"phone"`
	PasswordHash        string         `gorm:"not null" json:"-"`
	Role                UserRole       `gorm:"default:'manager'" json:"role"`
	FirstName           string         `json:"first_name"`
	LastName            string         `json:"last_name"`
	ShopID              *uuid.UUID     `gorm:"type:uuid" json:"shop_id"`
	IsPhoneVerified     bool           `gorm:"default:false" json:"is_phone_verified"`
	MustChangePassword  bool           `gorm:"default:false" json:"must_change_password"`
	VerificationCode    string         `json:"-"`
	VerificationExpires *time.Time     `json:"-"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	Account Account `gorm:"foreignKey:AccountID" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
