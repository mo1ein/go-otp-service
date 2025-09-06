// internal/handlers/health_handler.go
package handlers

import (
	"net/http"
	"otp-auth-service/internal/repository"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	userRepo repository.UserRepository
}

func NewHealthHandler(userRepo repository.UserRepository) *HealthHandler {
	return &HealthHandler{userRepo: userRepo}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Check database connection
	if err := h.userRepo.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "Database connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"version": "1.0.0",
	})
}
