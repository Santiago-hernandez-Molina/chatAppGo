package ports

import (
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type MessageRepo interface {
	GetMessagesByRoomId(roomId int) ([]models.MessageUser, error)
	SaveMessage(message *models.Message) error
}

type MessageUseCase interface {
	GetMessagesByRoomId(roomId int) ([]models.MessageUser, error)
	SaveMessage(message *models.Message) error
}

type MessageHandler interface {
	GetMessagesByRoomId(ctx *gin.Context)
}
