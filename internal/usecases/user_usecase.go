package usecases

import (
	"errors"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
)

type UserUseCase struct {
	repo            ports.UserRepo
	sessionManager  ports.SessionManager
	passwordManager ports.PasswordManager
	emailSender     ports.EmailSender
	userTask        ports.UserTask
}

func (useCase *UserUseCase) ActivateAccount(code int, email string) error {
	err := useCase.repo.ActivateAccount(code, email)
	return err
}

func (useCase *UserUseCase) DeleteInactiveUser(email string) error {
	err := useCase.repo.DeleteInactiveUser(email)
	return err
}

func (useCase *UserUseCase) DeleteUser(userId int) error {
	err := useCase.repo.DeleteUser(userId)
	return err
}

func (useCase *UserUseCase) GetCredentials(token string) (*models.User, error) {
	userClaims, err := useCase.GetCredentials(token)
	if err != nil {
		return nil, err
	}
	return userClaims, nil
}

func (useCase *UserUseCase) Login(user *models.User) (*models.UserWithToken, error) {
	userDB, err := useCase.repo.GetUserByEmail(user)
	if err != nil {
		return nil, err
	}
	success := useCase.passwordManager.ValidatePassword(
		userDB.Password,
		user.Password,
	)
	if !success {
		return nil, errors.New("Invalid Credentials")
	}
	tokenSTR, err := useCase.sessionManager.SignToken(userDB)
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

func (useCase *UserUseCase) Register(user *models.User) error {
	if len(user.Password) < 8 {
		return errors.New("Incorrect Password")
	}
	encryptedPassowrd, err := useCase.passwordManager.EncryptPassword(user.Password)
	if err != nil {
		return err
	}
    code, err := useCase.sessionManager.GenerateAuthCode()
	if err != nil {
		return err
	}

	user.Password = encryptedPassowrd
    user.Code = code

	err = useCase.repo.Register(user)
	if err != nil {
		return err
	}
	err = useCase.emailSender.SendRegisterConfirm(user, code)
	if err != nil {
        //To do: Delete User
		return err
	}
	go useCase.userTask.DeleteAccountTask(user.Email)

	return nil
}

var _ ports.UserUseCase = (*UserUseCase)(nil)

func NewUserUseCase(
	repo ports.UserRepo,
	sessionManager ports.SessionManager,
	passwordManager ports.PasswordManager,
	emailSender ports.EmailSender,
	userTask ports.UserTask,
) *UserUseCase {
	return &UserUseCase{
		repo:            repo,
		sessionManager:  sessionManager,
		passwordManager: passwordManager,
		emailSender:     emailSender,
		userTask:        userTask,
	}
}
