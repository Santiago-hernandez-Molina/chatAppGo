package services

import (
	"errors"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
)

type UserService struct {
	repo            ports.UserRepo
	sessionManager  ports.SessionManager
	passwordManager ports.PasswordManager
}

func (service *UserService) GetCredentials(token string) (*models.User, error) {
	userClaims, err := service.GetCredentials(token)
	if err != nil {
		return nil, err
	}
	return userClaims, nil
}

func (service *UserService) Login(user *models.User) (*models.UserWithToken, error) {
	userDB, err := service.repo.GetUserByEmail(user)
	if err != nil {
		return nil, err
	}
	success := service.passwordManager.ValidatePassword(
		userDB.Password,
		user.Password,
	)
	if !success {
		return nil, errors.New("Invalid Credentials")
	}
	tokenSTR, err := service.sessionManager.SignToken(userDB)
	if err != nil {
		return nil, err
	}
	userDB.Password = ""
	userWithToken := models.UserWithToken{
		User:  userDB,
		Token: tokenSTR,
	}
	return &userWithToken, nil
}

func (service *UserService) Register(user *models.User) error {
	if len(user.Password) < 8 {
		return errors.New("Incorrect Password")
	}
	encryptedPassowrd, err := service.passwordManager.EncryptPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = encryptedPassowrd
	err = service.repo.Register(user)
	if err != nil {
		return err
	}
	return nil
}

var _ ports.UserService = (*UserService)(nil)

func NewUserService(
    repo ports.UserRepo,
    sessionManager ports.SessionManager,
    passwordManager ports.PasswordManager,
) *UserService {
	return &UserService{
		repo:           repo,
		sessionManager: sessionManager,
        passwordManager: passwordManager,
	}
}
