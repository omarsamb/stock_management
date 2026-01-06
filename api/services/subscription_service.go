package services

import (
	"stock_management/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionService struct {
	DB *gorm.DB
}

func NewSubscriptionService(db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{DB: db}
}

// CreateSubscriptionRecord creates a new subscription record (pending)
func (s *SubscriptionService) CreateSubscriptionRecord(accountID uuid.UUID, amount float64, plan string) (*models.Subscription, error) {
	sub := &models.Subscription{
		ID:        uuid.New(),
		AccountID: accountID,
		Amount:    amount,
	}
	err := s.DB.Create(sub).Error
	return sub, err
}

// ConfirmPayment updates the account status based on the payment result
func (s *SubscriptionService) ConfirmPayment(accountID uuid.UUID, externalID string, amount float64) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var account models.Account
		if err := tx.First(&account, "id = ?", accountID).Error; err != nil {
			return err
		}

		// Update or create subscription record for history
		// In a real PayDunya integration, we'd look up the active transaction by externalID
		sub := models.Subscription{
			ID:            uuid.New(),
			AccountID:     accountID,
			ExternalPayID: externalID,
			Amount:        amount,
			ExpiryDate:    time.Now().AddDate(0, 1, 0), // Default 1 month
		}
		if err := tx.Create(&sub).Error; err != nil {
			return err
		}

		// Update Account Status
		newExpiry := time.Now().AddDate(0, 1, 0)
		account.Status = models.AccountStatusActive
		account.SubscriptionExpiresAt = &newExpiry

		return tx.Save(&account).Error
	})
}

// CheckAndSuspendAccounts flips accounts to read-only if expired
func (s *SubscriptionService) CheckAndSuspendAccounts() error {
	now := time.Now()
	// Find all active/trial accounts that have expired
	result := s.DB.Model(&models.Account{}).
		Where("(status = ? OR status = ?) AND subscription_expires_at < ?", models.AccountStatusActive, models.AccountStatusTrial, now).
		Update("status", models.AccountStatusReadOnly)

	return result.Error
}
