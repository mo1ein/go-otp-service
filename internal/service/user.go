package service

import (
	"otp-auth-service/internal/models"
	"otp-auth-service/internal/repository"
)

type UserService interface {
	GetUser(id uint) (*models.UserResponse, error)
	GetUsers(offset, limit int, search string) ([]models.UserResponse, int64, error)
	GetMe(userID uint) (*models.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUser(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (s *userService) GetUsers(offset, limit int, search string) ([]models.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAll(offset, limit, search)
	if err != nil {
		return nil, 0, err
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			ID:          user.ID,
			PhoneNumber: user.PhoneNumber,
			CreatedAt:   user.CreatedAt,
		})
	}

	return userResponses, total, nil
}

func (s *userService) GetMe(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
	}, nil
}
