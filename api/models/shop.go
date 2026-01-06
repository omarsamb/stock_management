package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shop struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	AccountID uuid.UUID      `gorm:"type:uuid;not null;index" json:"account_id"`
	Name      string         `gorm:"not null" json:"name"`
	Location  string         `json:"location"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Account Account `gorm:"foreignKey:AccountID" json:"-"`
}

func (s *Shop) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
