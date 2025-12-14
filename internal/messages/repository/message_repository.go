package repository

import (
	"Market_backend/internal/common"
	"Market_backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{db: common.DB}
}

func (r *MessageRepository) CreateMessage(msg *models.Message) error {
	msg.ID = uuid.New()
	msg.CreatedAt = time.Now()
	return r.db.Create(msg).Error
}

func (r *MessageRepository) GetAllMessages() ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Order("created_at DESC").Find(&messages).Error
	return messages, err
}
