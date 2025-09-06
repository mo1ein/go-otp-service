package model

import (
	"time"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PhoneNumber string    `json:"phone_number" gorm:"uniqueIndex"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserResponse struct {
	ID          uint      `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}
