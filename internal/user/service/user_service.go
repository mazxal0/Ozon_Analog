package service

import (
	"Market_backend/internal/user/repository"
	"Market_backend/models"

	"github.com/google/uuid"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.GetAllUsers()
}

func (s *UserService) GetMe(id uuid.UUID) (*models.User, error) {
	return s.repo.GetMe(id)
}
