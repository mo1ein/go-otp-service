// internal/handler/auth_handler.go
package handler

import (
	"net/http"
	"otp-auth-service/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// normalizePhoneNumber ensures phone number has + prefix
func normalizePhoneNumber(phone string) string {
	phone = strings.TrimSpace(phone)
	if !strings.HasPrefix(phone, "+") {
		phone = "+" + phone
	}
	return phone
}

type RequestOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
}

// RequestOTP godoc
// @Summary Request OTP for login/registration
// @Description Generate and send OTP to the provided phone number
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RequestOTPRequest true "Phone number"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/request-otp [post]
func (h *AuthHandler) RequestOTP(c *gin.Context) {
	var req RequestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.authService.RequestOTP(req.PhoneNumber)
	if err != nil {
		if err.Error() == "rate limit exceeded" {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// VerifyOTP godoc
// @Summary Verify OTP and login/register
// @Description Verify OTP and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "Phone number and OTP"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/verify-otp [post]
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.authService.VerifyOTP(req.PhoneNumber, req.OTP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
