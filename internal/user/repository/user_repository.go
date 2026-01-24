package repository

import (
	"Market_backend/internal/common"
	"Market_backend/internal/user/dto"
	"Market_backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: common.DB,
	}
}

// Можно добавить методы поиска по email, id, etc.
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", ID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetMe(userId uuid.UUID) (*models.User, error) {
	var user models.User

	err := r.db.First(&user, "users.id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) ChangeUser(
	userID uuid.UUID,
	change dto.UserChange,
) (*models.User, error) {

	// 1. Обновляем
	result := r.db.
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(change)

	if result.Error != nil {
		return nil, result.Error
	}

	// 2. Забираем обновлённого пользователя
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
