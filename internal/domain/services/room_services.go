package services

import (
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
)

type RoomService struct {
	repo ports.RoomRepo
}

func (service *RoomService) AddUserToRoom(userId int, roomId int) error {
    err := service.repo.AddUserToRoom(userId, roomId)
    if err != nil {
        return err
    }

    return nil
}


var _ ports.RoomService = (*RoomService)(nil)

func NewRoomService(repo ports.RoomRepo) *RoomService {
	return &RoomService{
		repo: repo,
	}
}

func (rs *RoomService) NewRoom(room models.Room, userId int) error {
	err := rs.repo.NewRoom(room, userId)
	if err != nil {
		return err
	}
	return nil
}

func (rs *RoomService) GetRoomById(roomId int) (*models.Room, error) {
	room, err := rs.repo.GetRoomById(roomId)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rs *RoomService) GetRoomsByUserId(userId int) ([]models.Room, error) {
	room, err := rs.repo.GetRoomsByUserId(userId)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rs *RoomService) GetUserRoom(userId int, roomId int) (*models.UserRoom, error) {
    userRoom, err := rs.repo.GetUserRoom(userId, roomId)
    if err != nil {
        return nil, err
    }
    return userRoom, nil
}
