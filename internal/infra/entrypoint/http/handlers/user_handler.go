package handlers

import (
	"net/http"
	"strconv"

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

func (handler *UserHandler) GetUsers(c *gin.Context) {
	usernameParam := c.Query("username")
	limitParam := c.Query("limit")
	offsetParam := c.Query("offset")
	authCookie, _ := c.Cookie("Authorization")
	claims, _ := handler.sessionManager.GetCredentials(authCookie)

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid limit param",
		})
		return
	}
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid offset param",
		})
		return
	}

	users, err := handler.userUseCase.GetUsersByUsername(claims.UserId, usernameParam, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot find the users",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

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
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Json incorrect data provided",
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Not found user with the given credentials",
			"err":     err.Error(),
		})
		return
	}
	c.SetSameSite(http.SameSiteDefaultMode)
	c.SetCookie(
		"Authorization",
		userWithToken.Token,
		14_400,
		"/",
		"",
		false,
		true,
	)
	c.JSON(http.StatusOK, &dtos.LoginResponse{
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
