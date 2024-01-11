package handlers

import (
	"net/http"
	"strconv"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/gin-gonic/gin"
)

type ContactRequestHandler struct {
	contactRequestUseCase ports.ContactRequestUseCase
	sessionManager        ports.SessionManager
}

func (handler *ContactRequestHandler) AcceptRequest(ctx *gin.Context) {
	authCookie, _ := ctx.Request.Cookie("Authorization")
	claims, _ := handler.sessionManager.GetCredentials(authCookie.Value)
	requestIDParam := ctx.Param("requestid")
	requestID, err := strconv.Atoi(requestIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid user param",
		})
		return
	}
	err = handler.contactRequestUseCase.AcceptRequest(
		requestID,
		claims.UserId,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Request Not found",
		})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{
		"message": "Request accepted",
	})
}

func (*ContactRequestHandler) GetReceivedRequests(ctx *gin.Context) {
	panic("unimplemented")
}

func (*ContactRequestHandler) GetSendedRequests(ctx *gin.Context) {
	panic("unimplemented")
}

func (handler *ContactRequestHandler) SendRequest(ctx *gin.Context) {
	authCookie, _ := ctx.Request.Cookie("Authorization")
	claims, _ := handler.sessionManager.GetCredentials(authCookie.Value)
	request := models.ContactRequest{}
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid user param",
		})
		return
	}
	request.FromUserId = claims.UserId
	err = handler.contactRequestUseCase.SendRequest(&request)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Cannot found the user",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Request send successfully",
	})
}

var _ ports.ContactRequestHandler = (*ContactRequestHandler)(nil)

func NewContactRequestHandler(
	contactRequestUseCase ports.ContactRequestUseCase,
	sessionManager ports.SessionManager,
) *ContactRequestHandler {
	return &ContactRequestHandler{
		contactRequestUseCase: contactRequestUseCase,
		sessionManager: sessionManager,
	}
}
