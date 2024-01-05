package handlers

import (
	"net/http"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/dtos"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService    ports.UserService
	sessionManager ports.SessionManager
}

var _ ports.UserHandler = (*UserHandler)(nil)

func (uh *UserHandler) Login(c *gin.Context) {
	request := dtos.LoginRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong reading json",
		})
		return
	}
	userWithToken, err := uh.UserService.Login(
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
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"Authorization",
		userWithToken.Token,
		14_400,
		"/",
        "netlify.app",
		true,
		true,
	)
	c.JSON(http.StatusAccepted, &dtos.LoginResponse{
		Id:       userWithToken.User.Id,
		Email:    userWithToken.User.Email,
		Username: userWithToken.User.Username,
	})
}

func (uh *UserHandler) Register(c *gin.Context) {
	request := dtos.RegisterRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error reading data",
		})
	}
	err = uh.UserService.Register(&models.User{
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
		"message": "User Created Succesfully",
	})
}

func NewUserHandler(sessionManager ports.SessionManager, userService ports.UserService) *UserHandler {
	return &UserHandler{
		sessionManager: sessionManager,
		UserService:    userService,
	}
}
