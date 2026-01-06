package services

import (
	"fmt"
	"math/rand"
	"time"
)

type WhatsAppService struct{}

func NewWhatsAppService() *WhatsAppService {
	return &WhatsAppService{}
}

// SendOTP simulates sending a verification code via WhatsApp.
func (s *WhatsAppService) SendOTP(phone, code string) error {
	// In a real implementation, you would call a WhatsApp API here (Twilio, UltraMsg, etc.)
	fmt.Printf("\n--- WHATSAPP OTP [%s] sent to %s ---\n\n", code, phone)
	return nil
}

// GenerateOTP creates a random 6-digit code.
func (s *WhatsAppService) GenerateOTP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}
