package repository

import (
	"context"
	"otp-auth-service/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OTPRepository interface {
	StoreOTP(phoneNumber, otp string, expiration time.Duration) error
	GetOTP(phoneNumber string) (string, error)
	IncrementRequestCount(phoneNumber string, expiration time.Duration) (int, error)
	RecordOTPRequest(phoneNumber string, successful bool) error
	GetRequestCount(phoneNumber string, since time.Time) (int, error)
	GetSuccessfulRequestCount(phoneNumber string, since time.Time) (int, error)
}

type otpRepository struct {
	client *redis.Client
	db     *gorm.DB
}

func NewOTPRepository(client *redis.Client, db *gorm.DB) OTPRepository {
	return &otpRepository{client: client, db: db}
}

func (r *otpRepository) StoreOTP(phoneNumber, otp string, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, "otp:"+phoneNumber, otp, expiration).Err()
}

func (r *otpRepository) GetOTP(phoneNumber string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, "otp:"+phoneNumber).Result()
}

func (r *otpRepository) IncrementRequestCount(phoneNumber string, expiration time.Duration) (int, error) {
	// Use database for rate limiting instead of Redis
	since := time.Now().UTC().Add(-expiration)
	return r.GetRequestCount(phoneNumber, since)
}

func (r *otpRepository) RecordOTPRequest(phoneNumber string, successful bool) error {
	otpRequest := &model.OTPRequest{
		PhoneNumber: phoneNumber,
		RequestedAt: time.Now().UTC(),
		Successful:  successful,
	}
	return r.db.Create(otpRequest).Error
}

func (r *otpRepository) GetRequestCount(phoneNumber string, since time.Time) (int, error) {
	var count int64
	err := r.db.Model(&model.OTPRequest{}).
		Where("phone_number = ? AND requested_at >= ?", phoneNumber, since.UTC()).
		Count(&count).Error

	return int(count), err
}

func (r *otpRepository) GetSuccessfulRequestCount(phoneNumber string, since time.Time) (int, error) {
	var count int64
	err := r.db.Model(&model.OTPRequest{}).
		Where("phone_number = ? AND requested_at >= ? AND successful = ?", phoneNumber, since, true).
		Count(&count).Error
	return int(count), err
}
