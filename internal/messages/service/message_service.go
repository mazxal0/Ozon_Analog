package service

import (
	"Market_backend/internal/messages/repository"
	"Market_backend/models"
)

type MessageService struct {
	repo *repository.MessageRepository
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) SendMessage(name, email, phone, text string) error {
	msg := &models.Message{
		Name:  name,
		Email: email,
		Phone: phone,
		Text:  text,
	}
	return s.repo.CreateMessage(msg)
}

func (s *MessageService) GetMessages() ([]models.Message, error) {
	return s.repo.GetAllMessages()
}
