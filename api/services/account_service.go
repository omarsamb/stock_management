package services

import (
	"errors"
	"stock_management/models"
	"stock_management/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountService struct {
	DB        *gorm.DB
	JWTSecret string
}

func NewAccountService(db *gorm.DB, jwtSecret string) *AccountService {
	return &AccountService{DB: db, JWTSecret: jwtSecret}
}

// Register creates an Account and an Owner User in one transaction.
func (s *AccountService) Register(phone, password, companyName string) (*models.User, *models.Account, error) {
	var user *models.User
	var account *models.Account

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Create Account
		account = &models.Account{
			CompanyName:      companyName,
			SubscriptionPlan: "basic",
			Status:           models.AccountStatusTrial,
		}
		if err := tx.Create(account).Error; err != nil {
			return err
		}

		// 2. Hash Password
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return err
		}

		// 3. Create Owner User
		user = &models.User{
			AccountID:    account.ID,
			Phone:        phone,
			PasswordHash: hashedPassword,
			Role:         models.RoleOwner,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return nil
	})

	return user, account, err
}

// Authenticate verifies credentials and returns the user if successful.
func (s *AccountService) Authenticate(phone, password string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	match, err := utils.ComparePasswords(password, user.PasswordHash)
	if err != nil || !match {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

// UpdateProfile updates user personal information.
func (s *AccountService) UpdateProfile(userID uuid.UUID, firstName, lastName string) error {
	return s.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
	}).Error
}

// UpdatePlan updates the account's subscription plan.
func (s *AccountService) UpdatePlan(accountID uuid.UUID, plan string) error {
	return s.DB.Model(&models.Account{}).Where("id = ?", accountID).Update("subscription_plan", plan).Error
}

func (s *AccountService) InviteUser(accountID uuid.UUID, phone, role string) (*models.User, error) {
	user := &models.User{
		AccountID: accountID,
		Phone:     phone,
		Role:      models.UserRole(role),
	}
	err := s.DB.Create(user).Error
	return user, err
}
