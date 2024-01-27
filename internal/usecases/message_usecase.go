package usecases

import (
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
)

type MessageUseCase struct {
	repo ports.MessageRepo
}

func (useCase *MessageUseCase) SaveMessage(message *models.Message) (int, error) {
	id, err := useCase.repo.SaveMessage(message)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (useCase *MessageUseCase) GetMessagesByRoomId(roomId int) ([]models.MessageUser, error) {
	messages, err := useCase.repo.GetMessagesByRoomId(roomId)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

var _ ports.MessageUseCase = (*MessageUseCase)(nil)

func NewMessageUseCase(repo ports.MessageRepo) *MessageUseCase {
	return &MessageUseCase{
		repo: repo,
	}
}
