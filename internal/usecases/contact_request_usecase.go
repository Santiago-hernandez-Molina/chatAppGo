package usecases

import (
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
)

type ContactRequestUseCase struct {
	repo     ports.ContactRequestRepo
	roomRepo ports.RoomRepo
}

func (useCase *ContactRequestUseCase) AcceptRequest(requestId int, userId int) error {
	request, err := useCase.repo.GetRequestByToUserId(requestId, userId)
	if err != nil {
		return err
	}
	room := models.Room{
		Name: "Contact chat",
		Users: []models.UserRoom{
			{UserId: request.FromUserId, Role: "user"},
			{UserId: request.ToUserId, Role: "user"},
		},
	}
	room.Type = models.RoomType(models.Contact)
	err = useCase.roomRepo.NewRoom(&room)
	return err
}

func (useCase *ContactRequestUseCase) GetReceivedRequests(userid int) ([]models.ContactRequest, error) {
	requests, err := useCase.repo.GetReceivedRequests(userid)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (useCase *ContactRequestUseCase) GetSendedRequests(userid int) ([]models.ContactRequest, error) {
	requests, err := useCase.repo.GetSendedRequests(userid)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (useCase *ContactRequestUseCase) SendRequest(request *models.ContactRequest) error {
	err := useCase.repo.SaveRequest(request)
	return err
}

var _ ports.ContactRequestUseCase = (*ContactRequestUseCase)(nil)

func NewContactRequestUseCase(
	contactRequestRepo ports.ContactRequestRepo,
	roomRepo ports.RoomRepo,
) *ContactRequestUseCase {
	return &ContactRequestUseCase{
		repo:     contactRequestRepo,
		roomRepo: roomRepo,
	}
}
