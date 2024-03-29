package usecases

import (
	"errors"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/exceptions"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
)

type ContactRequestUseCase struct {
	repo     ports.ContactRequestRepo
	userRepo ports.UserRepo
	roomRepo ports.RoomRepo
}

func (useCase *ContactRequestUseCase) AcceptRequest(requestId int, userId int) error {
	request, err := useCase.repo.GetRequestById(requestId)
	if err != nil {
		return err
	}
	if request.Accepted == true {
		return errors.New("request already accepted")
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
	if err != nil {
		return err
	}
	err = useCase.repo.UpdateRequestStatus(true, requestId)
	return err
}

func (useCase *ContactRequestUseCase) GetReceivedRequests(userid int) ([]models.ContactRequestWithUser, error) {
	requests, err := useCase.repo.GetReceivedRequests(userid)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (useCase *ContactRequestUseCase) GetSendedRequests(userid int) ([]models.ContactRequestWithUser, error) {
	requests, err := useCase.repo.GetSendedRequests(userid)
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (useCase *ContactRequestUseCase) SendRequest(request *models.ContactRequest) error {
	_, err := useCase.userRepo.GetUserById(request.ToUserId)
	if err != nil {
		return &exceptions.UserNotFound{}
	}
	_, err = useCase.repo.GetRequestByToUserId(request.ToUserId, request.FromUserId)
	if err == nil {
		return errors.New("request already sended")
	}
	err = useCase.repo.SaveRequest(request)
	return err
}

var _ ports.ContactRequestUseCase = (*ContactRequestUseCase)(nil)

func NewContactRequestUseCase(
	contactRequestRepo ports.ContactRequestRepo,
	roomRepo ports.RoomRepo,
	userRepo ports.UserRepo,
) *ContactRequestUseCase {
	return &ContactRequestUseCase{
		repo:     contactRequestRepo,
		roomRepo: roomRepo,
		userRepo: userRepo,
	}
}
