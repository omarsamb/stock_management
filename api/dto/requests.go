package dto

import "github.com/google/uuid"

type CreateArticleRequest struct {
	Name         string     `json:"name" binding:"required"`
	SKU          string     `json:"sku"`
	Description  string     `json:"description"`
	CategoryID   *uuid.UUID `json:"category_id"`
	BrandID      *uuid.UUID `json:"brand_id"`
	MinThreshold int        `json:"min_threshold"`
	Price        float64    `json:"price"`
}

type UpdateArticleRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description"`
	MinThreshold int     `json:"min_threshold"`
	Price        float64 `json:"price"`
}

type RecordMovementRequest struct {
	ShopID    uuid.UUID `json:"shop_id" binding:"required"`
	ArticleID uuid.UUID `json:"article_id" binding:"required"`
	Type      string    `json:"type" binding:"required"` // in, out, adjust
	Qty       int       `json:"qty" binding:"required"`
	Reason    string    `json:"reason"`
	DeviceID  string    `json:"device_id"`
}

type TransferStockRequest struct {
	FromShopID uuid.UUID `json:"from_shop_id" binding:"required"`
	ToShopID   uuid.UUID `json:"to_shop_id" binding:"required"`
	ArticleID  uuid.UUID `json:"article_id" binding:"required"`
	Qty        int       `json:"qty" binding:"required"`
	Reason     string    `json:"reason"`
	DeviceID   string    `json:"device_id"`
}

type RegisterRequest struct {
	Phone       string `json:"phone" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	CompanyName string `json:"company_name" binding:"required"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type SelectPlanRequest struct {
	Plan string `json:"plan" binding:"required"`
}

type CreateShopRequest struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location"`
}

type InviteUserRequest struct {
	FirstName string     `json:"first_name" binding:"required"`
	LastName  string     `json:"last_name" binding:"required"`
	Phone     string     `json:"phone" binding:"required"`
	Role      string     `json:"role" binding:"required"`
	ShopID    *uuid.UUID `json:"shop_id"`
}

type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
type UpdateUserRequest struct {
	FirstName string     `json:"first_name" binding:"required"`
	LastName  string     `json:"last_name" binding:"required"`
	Phone     string     `json:"phone" binding:"required"`
	Role      string     `json:"role" binding:"required"`
	ShopID    *uuid.UUID `json:"shop_id"`
}

type UpdateThemeRequest struct {
	PrimaryColor    string `json:"primary_color"`
	BackgroundImage string `json:"background_image"`
}
