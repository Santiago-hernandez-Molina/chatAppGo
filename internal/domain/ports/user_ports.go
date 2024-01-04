package ports

import (
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type UserRepo interface {
	GetUserByEmail(user *models.User) (*models.User, error)
	Register(user *models.User) error
}

type UserService interface {
	Login(user *models.User) (*models.UserWithToken, error)
	Register(user *models.User) error
	GetCredentials(token string) (*models.User, error)
}

type UserHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type SessionManager interface {
	GetCredentials(token string) (*models.UserClaims, error)
	SignToken(user *models.User) (string, error)
}

type PasswordManager interface {
	ValidatePassword(encryptedPassword string, requestedPassword string) bool
	EncryptPassword(passsword string) (string, error)
}
