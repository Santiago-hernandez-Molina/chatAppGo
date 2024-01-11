package ports

import (
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type RoomRepo interface {
	GetRoomById(roomId int) (*models.Room, error)
	GetRoomsByUserId(userId int) ([]models.Room, error)
	NewRoom(room *models.Room) error
	AddUserToRoom(userId int, roomId int) error
	GetUserRoom(roomId int, userId int) (*models.UserRoom, error)
}

type RoomUseCase interface {
	GetRoomById(roomId int, userId int) (*models.Room, error)
	GetRoomsByUserId(userId int) ([]models.Room, error)
	NewRoom(room models.Room, userId int) error
	AddUserToRoom(userId int, roomId int) error
	GetUserRoom(userId int, roomId int) (*models.UserRoom, error)
}

type RoomHandler interface {
	GetRoomsByUserId(ctx *gin.Context)
	GetRoomById(ctx *gin.Context)
	NewRoom(ctx *gin.Context)
	AddUserToRoom(ctx *gin.Context)
	ConnectToRoom(ctx *gin.Context)
}
