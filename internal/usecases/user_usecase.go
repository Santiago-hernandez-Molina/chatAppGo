package usecases

import (
	"errors"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/exceptions"
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

func (useCase *UserUseCase) GetUsersByUsername(
	userId int,
	filter string,
	size int,
	offset int,
) (*models.PaginatedModel[models.UserContact], error) {
	count, err := useCase.repo.GetUsersCount(filter)
	if err != nil {
		return nil, err
	}
	if offset > count || offset < 0 || size < 0 {
		return nil, errors.New("invalid offset")
	}
	if size == 0 {
		size = 5
	}

	paginatedUsers, err := useCase.repo.GetUsersByUsername(
		userId,
		filter,
		size,
		offset,
	)
	if err != nil {
		return nil, err
	}

	paginatedUsers.Count = count
	return paginatedUsers, nil
}

func (useCase *UserUseCase) ActivateAccount(code int, email string) error {
	err := useCase.repo.ActivateAccount(code, email)
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
		return nil, errors.New("something went wrong")
	}
	if userDB.Status == false {
		return nil, errors.New("The account has not been actived")
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
	_, err := useCase.repo.GetUserByEmail(user)
	if err == nil {
		return &exceptions.DuplicatedUser{}
	}
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
		user, errUser := useCase.repo.GetUserByEmail(user)
		if errUser != nil {
			return err
		}
		useCase.repo.DeleteUser(user.Id)
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
