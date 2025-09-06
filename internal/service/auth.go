package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"otp-auth-service/internal/models"
	"otp-auth-service/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	RequestOTP(phoneNumber string) error
	VerifyOTP(phoneNumber, otp string) (string, error)
	GenerateJWT(user *models.User) (string, error)
}

type authService struct {
	userRepo   repository.UserRepository
	otpRepo    repository.OTPRepository
	jwtSecret  string
	otpExpiry  time.Duration
	rateLimit  int
	rateWindow time.Duration
}

func NewAuthService(userRepo repository.UserRepository, otpRepo repository.OTPRepository) AuthService {
	return &authService{
		userRepo:   userRepo,
		otpRepo:    otpRepo,
		otpExpiry:  2 * time.Minute,
		rateLimit:  3,
		rateWindow: 10 * time.Minute,
	}
}

func (s *authService) RequestOTP(phoneNumber string) error {
	// Check rate limiting
	count, err := s.otpRepo.IncrementRequestCount(phoneNumber, s.rateWindow)
	if err != nil {
		return err
	}

	if count > s.rateLimit {
		// Record failed request due to rate limiting
		s.otpRepo.RecordOTPRequest(phoneNumber, false)
		return fmt.Errorf("rate limit exceeded")
	}

	// Generate OTP
	otp, err := generateOTP(6)
	if err != nil {
		// Record failed request due to OTP generation error
		s.otpRepo.RecordOTPRequest(phoneNumber, false)
		return err
	}

	// Store OTP
	err = s.otpRepo.StoreOTP(phoneNumber, otp, s.otpExpiry)
	if err != nil {
		// Record failed request due to storage error
		s.otpRepo.RecordOTPRequest(phoneNumber, false)
		return err
	}

	// Record successful request
	s.otpRepo.RecordOTPRequest(phoneNumber, true)

	// Print OTP to console (instead of sending SMS)
	fmt.Printf("OTP for %s: %s\n", phoneNumber, otp)

	return nil
}

func (s *authService) VerifyOTP(phoneNumber, otp string) (string, error) {
	// Get stored OTP
	storedOTP, err := s.otpRepo.GetOTP(phoneNumber)
	if err != nil {
		return "", fmt.Errorf("invalid or expired OTP")
	}

	// Verify OTP
	if storedOTP != otp {
		return "", fmt.Errorf("invalid OTP")
	}

	// Find or create user
	user, err := s.userRepo.FindByPhoneNumber(phoneNumber)
	if err != nil {
		// User doesn't exist, create new one
		user = &models.User{
			PhoneNumber: phoneNumber,
			CreatedAt:   time.Now().UTC(),
		}
		err = s.userRepo.Create(user)
		if err != nil {
			return "", err
		}
	}

	// Generate JWT token
	token, err := s.GenerateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"phone":   user.PhoneNumber,
		"exp":     time.Now().UTC().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func generateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}
