package handlers

import (
	"net/http"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/dtos"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase    ports.UserUseCase
	userTask       ports.UserTask
	sessionManager ports.SessionManager
}

var _ ports.UserHandler = (*UserHandler)(nil)

func (handler *UserHandler) ActivateAccount(c *gin.Context) {
	var activationRequest dtos.ActivationRequest
	err := c.BindJSON(&activationRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}
	err = handler.userUseCase.ActivateAccount(
		activationRequest.Code,
		activationRequest.Email,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot active the account, try register again",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Account activated successfully",
	})
}

func (handler *UserHandler) Login(c *gin.Context) {
	request := dtos.LoginRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong reading json",
		})
		return
	}
	userWithToken, err := handler.userUseCase.Login(
		&models.User{
			Email:    request.Email,
			Password: request.Password,
		},
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found user with the given credentials",
			"err":     err.Error(),
		})
		return
	}
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(
		"Authorization",
		userWithToken.Token,
		14_400,
		"/",
		"",
		true,
		true,
	)
	c.JSON(http.StatusAccepted, &dtos.LoginResponse{
		Id:       userWithToken.User.Id,
		Email:    userWithToken.User.Email,
		Username: userWithToken.User.Username,
	})
}

func (handler *UserHandler) Register(c *gin.Context) {
	request := dtos.RegisterRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error reading data",
		})
	}
	err = handler.userUseCase.Register(&models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Check your email",
	})
}

func NewUserHandler(
	sessionManager ports.SessionManager,
	userUseCase ports.UserUseCase,
	userTask ports.UserTask,
) *UserHandler {
	return &UserHandler{
		userTask:       userTask,
		sessionManager: sessionManager,
		userUseCase:    userUseCase,
	}
}
