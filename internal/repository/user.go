package repository

import (
	"otp-auth-service/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByPhoneNumber(phoneNumber string) (*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindAll(offset, limit int, search string) ([]model.User, int64, error)
	HealthCheck() error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByPhoneNumber(phoneNumber string) (*model.User, error) {
	var user model.User
	err := r.db.Where("phone_number = ?", phoneNumber).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindAll(offset, limit int, search string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	if search != "" {
		query = query.Where("phone_number ILIKE ?", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *userRepository) HealthCheck() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
