package models

import "time"

type OTPRequest struct {
	ID          uint      `gorm:"primaryKey"`
	PhoneNumber string    `gorm:"column:phone_number"`
	RequestedAt time.Time `gorm:"column:requested_at"`
	Successful  bool      `gorm:"column:successful"`
}

func (*OTPRequest) TableName() string {
	return "otp_requests"
}

type OTPRequestResponse struct {
	ID          uint      `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	RequestedAt time.Time `json:"requested_at"`
	Successful  bool      `json:"successful"`
}
