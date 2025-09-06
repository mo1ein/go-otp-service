package handler

import (
	"net/http"
	"otp-auth-service/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type OTPStatsHandler struct {
	otpRepo repository.OTPRepository
}

func NewOTPStatsHandler(otpRepo repository.OTPRepository) *OTPStatsHandler {
	return &OTPStatsHandler{otpRepo: otpRepo}
}

// GetOTPStats godoc
// @Summary Get OTP request statistics
// @Description Get statistics about OTP requests for a phone number
// @Tags otp
// @Accept json
// @Produce json
// @Param phone query string true "Phone number"
// @Param hours query int false "Number of hours to look back" default(24)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /otp/stats [get]
func (h *OTPStatsHandler) GetOTPStats(c *gin.Context) {
	phoneNumber := strings.TrimSpace(c.Query("phone"))
	if phoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	// Ensure phone number has + prefix
	if !strings.HasPrefix(phoneNumber, "+") {
		phoneNumber = "+" + phoneNumber
	}

	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hours parameter"})
		return
	}

	since := time.Now().UTC().Add(-time.Duration(hours) * time.Hour)

	// Get total request count
	totalCount, err := h.otpRepo.GetRequestCount(phoneNumber, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get request count"})
		return
	}

	// Get successful request count
	successfulCount, err := h.otpRepo.GetSuccessfulRequestCount(phoneNumber, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get successful request count"})
		return
	}

	failedCount := totalCount - successfulCount
	successRate := 0.0
	if totalCount > 0 {
		successRate = float64(successfulCount) / float64(totalCount) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"phone_number":        phoneNumber,
		"hours_looked_back":   hours,
		"total_requests":      totalCount,
		"successful_requests": successfulCount,
		"failed_requests":     failedCount,
		"success_rate":        successRate,
		"since":               since.Format(time.RFC3339),
		"timezone":            "UTC",
	})
}
