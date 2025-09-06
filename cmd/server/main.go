package main

import (
	"context"
	"fmt"
	"log"
	"otp-auth-service/internal/config"
	"otp-auth-service/internal/handlers"
	"otp-auth-service/internal/middleware"
	"otp-auth-service/internal/models"
	"otp-auth-service/internal/repository"
	"otp-auth-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database
	// todo: fix ssl mode
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Database.Postgres.User, cfg.Database.Postgres.Password, cfg.Database.Postgres.Host, cfg.Database.Postgres.Port, cfg.Database.Postgres.DatabaseName)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if cfg.Database.Redis.Port == "" {
		cfg.Database.Redis.Port = "6379"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Database.Redis.Host, cfg.Database.Redis.Port),
		Password: cfg.Database.Redis.Password,
		DB:       cfg.Database.Redis.Database,
	})

	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate models
	db.AutoMigrate(&models.User{}, &models.OTPRequest{})

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(redisClient, db)

	// Initialize services
	authService := service.NewAuthService(userRepo, otpRepo)
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	otpStatsHandler := handlers.NewOTPStatsHandler(otpRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware()

	// Setup router
	router := gin.Default()

	// Auth routes
	router.POST("/auth/request-otp", authHandler.RequestOTP)
	router.POST("/auth/verify-otp", authHandler.VerifyOTP)

	// Protected routes
	router.GET("/me", authMiddleware.ValidateToken, userHandler.GetMe)

	// User routes (protected)
	userRoutes := router.Group("/users")
	userRoutes.Use(authMiddleware.ValidateToken)
	{
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.GET("/", userHandler.GetUsers)
	}

	// OTP stats route (public for monitoring)
	router.GET("/otp/stats", otpStatsHandler.GetOTPStats)

	// Start server
	router.Run(fmt.Sprintf("%s:%d", cfg.HTTP.APIHost, cfg.HTTP.APIPort))
}
