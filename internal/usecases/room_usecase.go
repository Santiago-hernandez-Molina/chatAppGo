package usecases

import (
	"errors"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/exceptions"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
)

type RoomUseCase struct {
	repo     ports.RoomRepo
	userRepo ports.UserRepo
}

func (useCase *RoomUseCase) AddUserToRoom(userId int, roomId int) error {
	_, err := useCase.repo.GetUserRoom(userId, roomId)
	if err == nil {
		return &exceptions.DuplicatedUser{}
	}
	if !errors.Is(err, &exceptions.UserNotFound{}) {
		return &exceptions.AccesDataException{}
	}
	_, err = useCase.userRepo.GetUserById(userId)
	if err != nil {
		return err
	}
	err = useCase.repo.AddUserToRoom(userId, roomId)
	if err != nil {
		return err
	}

	return nil
}

func (useCase *RoomUseCase) NewRoom(room models.Room, userId int) error {
	room.Type = models.RoomType(models.Group)
	room.Users = []models.UserRoom{
		{UserId: userId, RoomId: room.Id, Role: "admin"},
	}
	err := useCase.repo.NewRoom(&room)
	if err != nil {
		return err
	}
	return nil
}

func (useCase *RoomUseCase) GetRoomById(roomId int, userId int) (
	*models.Room,
	error,
) {
	room, err := useCase.repo.GetRoomById(roomId)
	if err != nil {
		return nil, err
	}
	if room.Type != models.RoomType(models.Contact) {
		return room, nil
	}
	if room.Users[0].UserId == userId {
		user2, err := useCase.userRepo.GetUserById(room.Users[1].UserId)
		if err != nil {
			return nil, &exceptions.UserNotFound{}
		}
		room.Name = user2.Username
	} else {
		user2, err := useCase.userRepo.GetUserById(room.Users[0].UserId)
		if err != nil {
			return nil, &exceptions.UserNotFound{}
		}
		room.Name = user2.Username
	}
	return room, nil
}

func (useCase *RoomUseCase) GetRoomsByUserId(userId int) (
	[]models.Room,
	error,
) {
	rooms, err := useCase.repo.GetRoomsByUserId(userId)
	if err != nil {
		return nil, err
	}
	for i, room := range rooms {
		if room.Type == models.RoomType(models.Group) {
			continue
		}
		if room.Users[0].UserId == userId {
			user2, err := useCase.userRepo.GetUserById(room.Users[1].UserId)
			if err != nil {
				rooms[i].Name = "Not Found"
				continue
			}
			rooms[i].Name = user2.Username
		} else {
			user2, err := useCase.userRepo.GetUserById(room.Users[0].UserId)
			if err != nil {
				rooms[i].Name = "Not Found"
				continue
			}
			rooms[i].Name = user2.Username
		}
	}
	return rooms, nil
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

func NewRoomUseCase(repo ports.RoomRepo, userRepo ports.UserRepo) *RoomUseCase {
	return &RoomUseCase{
		repo:     repo,
		userRepo: userRepo,
	}
}
