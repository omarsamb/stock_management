package handlers

import (
	"net/http"
	"stock_management/dto"
	"stock_management/models"
	"stock_management/services"
	"stock_management/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	Service         *services.AccountService
	WhatsAppService *services.WhatsAppService
	JWTSecret       string
}

func NewAuthHandler(s *services.AccountService, w *services.WhatsAppService, jwtSecret string) *AuthHandler {
	return &AuthHandler{Service: s, WhatsAppService: w, JWTSecret: jwtSecret}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _, err := h.Service.Register(req.Phone, req.Password, req.CompanyName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate and send OTP
	code := h.WhatsAppService.GenerateOTP()
	expires := time.Now().Add(10 * time.Minute)

	// Update user with verification code
	h.Service.DB.Model(user).Updates(map[string]interface{}{
		"verification_code":    code,
		"verification_expires": &expires,
	})

	h.WhatsAppService.SendOTP(user.Phone, code)

	c.JSON(http.StatusCreated, gin.H{
		"message": "compte créé. veuillez vérifier votre téléphone via le code envoyé sur WhatsApp.",
		"phone":   user.Phone,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.Authenticate(req.Phone, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !user.IsPhoneVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "votre numéro n'est pas vérifié", "requires_verification": true})
		return
	}

	token, err := utils.GenerateToken(user, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":                 user,
		"account":              user.Account,
		"token":                token,
		"must_change_password": user.MustChangePassword,
	})
}

func (h *AuthHandler) VerifyPhone(c *gin.Context) {
	var req dto.VerifyPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.Service.DB.Where("phone = ?", req.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "utilisateur non trouvé"})
		return
	}

	if user.VerificationCode != req.Code || user.VerificationExpires.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code invalide ou expiré"})
		return
	}

	h.Service.DB.Model(&user).Updates(map[string]interface{}{
		"is_phone_verified":    true,
		"verification_code":    "",
		"verification_expires": nil,
	})

	token, _ := utils.GenerateToken(&user, h.JWTSecret)

	c.JSON(http.StatusOK, gin.H{
		"message": "téléphone vérifié avec succès",
		"token":   token,
		"user":    user,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in context"})
		return
	}

	if err := h.Service.UpdateProfile(userID, req.FirstName, req.LastName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}

func (h *AuthHandler) InviteUser(c *gin.Context) {
	var req dto.InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if requester is Owner
	role := c.GetString("role")
	if role != string(models.RoleOwner) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only owners can invite users"})
		return
	}

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	user, tempPassword, err := h.Service.InviteUser(accountID, req.Phone, req.Role, req.FirstName, req.LastName, req.ShopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "user invited successfully",
		"user":          user,
		"temp_password": tempPassword,
	})
}

func (h *AuthHandler) ListUsers(c *gin.Context) {
	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	users, err := h.Service.GetUsersByAccount(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.Param("id")
	userID, _ := uuid.Parse(userIDStr)

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	if err := h.Service.UpdateUser(accountID, userID, req.FirstName, req.LastName, req.Phone, req.Role, req.ShopID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *AuthHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, _ := uuid.Parse(userIDStr)

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	if err := h.Service.DeleteUser(accountID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	if err := h.Service.ChangePassword(userID, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

func (h *AuthHandler) UpdateTheme(c *gin.Context) {
	var req dto.UpdateThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accountIDStr := c.GetString("account_id")
	accountID, _ := uuid.Parse(accountIDStr)

	if err := h.Service.UpdateAccountTheme(accountID, req.PrimaryColor, req.BackgroundImage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "theme updated successfully"})
}
