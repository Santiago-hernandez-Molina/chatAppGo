package ports

import (
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type ContactRequestRepo interface {
	SaveRequest(request *models.ContactRequest) error
	GetSendedRequests(userid int) ([]models.ContactRequest, error)
	GetReceivedRequests(userid int) ([]models.ContactRequest, error)
	GetRequestByToUserId(userId int) (*models.ContactRequest, error)
	GetRequestById(requestId int) (*models.ContactRequest, error)
}

type ContactRequestUseCase interface {
	SendRequest(request *models.ContactRequest) error
	AcceptRequest(requestId int, userId int) error
	GetSendedRequests(userid int) ([]models.ContactRequest, error)
	GetReceivedRequests(userid int) ([]models.ContactRequest, error)
}

type ContactRequestHandler interface {
	SendRequest(ctx *gin.Context)
	AcceptRequest(ctx *gin.Context)
	GetSendedRequests(ctx *gin.Context)
	GetReceivedRequests(ctx *gin.Context)
}
