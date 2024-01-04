package services

import (
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
)

type MessageService struct {
	repo ports.MessageRepo
}

func (ms *MessageService) SaveMessage(message *models.Message) error {
	err := ms.repo.SaveMessage(message)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MessageService) GetMessagesByRoomId(roomId int) ([]models.MessageUser, error) {
	messages, err := ms.repo.GetMessagesByRoomId(roomId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

var _ ports.MessageService = (*MessageService)(nil)

func NewMessageService(repo ports.MessageRepo) *MessageService {
	return &MessageService{
		repo: repo,
	}
}
