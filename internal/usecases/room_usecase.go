package usecases

import (
	"errors"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/exceptions"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
)

type RoomUseCase struct {
	repo ports.RoomRepo
}

func (useCase *RoomUseCase) AddUserToRoom(userId int, roomId int) error {
	_, err := useCase.repo.GetUserRoom(userId, roomId)
	if err == nil {
		return &exceptions.DuplicatedUser{}
	}
	if !errors.Is(err, &exceptions.UserNotFound{}) {
		return &exceptions.AccesDataException{}
	}
	err = useCase.repo.AddUserToRoom(userId, roomId)
	if err != nil {
		return err
	}

	return nil
}

func (useCase *RoomUseCase) NewRoom(room models.Room, userId int) error {
	err := useCase.repo.NewRoom(room, userId)
	if err != nil {
		return err
	}
	return nil
}

func (useCase *RoomUseCase) GetRoomById(roomId int) (
	*models.Room,
	error,
) {
	room, err := useCase.repo.GetRoomById(roomId)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (useCase *RoomUseCase) GetRoomsByUserId(userId int) (
	[]models.Room,
	error,
) {
	room, err := useCase.repo.GetRoomsByUserId(userId)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (useCase *RoomUseCase) GetUserRoom(userId int, roomId int) (
	*models.UserRoom,
	error,
) {
	userRoom, err := useCase.repo.GetUserRoom(
		userId,
		roomId,
	)
	if err != nil {
		return nil, err
	}
	return userRoom, nil
}

var _ ports.RoomUseCase = (*RoomUseCase)(nil)

func NewRoomUseCase(repo ports.RoomRepo) *RoomUseCase {
	return &RoomUseCase{
		repo: repo,
	}
}
