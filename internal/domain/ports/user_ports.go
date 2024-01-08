package ports

import (
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/gin-gonic/gin"
)

type UserRepo interface {
	GetUserByEmail(user *models.User) (*models.User, error)
	Register(user *models.User) error
	DeleteUser(userId int) error
	DeleteInactiveUser(email string) error
    ActivateAccount(code int, email string) error
}

type UserUseCase interface {
	Login(user *models.User) (*models.UserWithToken, error)
	Register(user *models.User) error
	DeleteUser(userId int) error
	GetCredentials(token string) (*models.User, error)
    ActivateAccount(code int, email string) error
}

type UserTask interface {
	DeleteAccountTask(email string) error
}

type UserHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
    ActivateAccount(c *gin.Context)
}

// Session Related

type SessionManager interface {
	GetCredentials(token string) (*models.UserClaims, error)
	SignToken(user *models.User) (string, error)
	GenerateAuthCode() (int, error)
}

type PasswordManager interface {
	ValidatePassword(encryptedPassword string, requestedPassword string) bool
	EncryptPassword(passsword string) (string, error)
}
