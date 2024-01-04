package authentication

import (
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"golang.org/x/crypto/bcrypt"
)

type PasswordManager struct{}

func (*PasswordManager) EncryptPassword(password string) (string, error) {
	hashP, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hashP), nil
}

func (manager *PasswordManager) ValidatePassword(
	encryptedPassword string,
	requestedPassword string,
) bool {
	if err := bcrypt.CompareHashAndPassword(
		[]byte(encryptedPassword),
		[]byte(requestedPassword),
	); err != nil {
		return false
	}
	return true
}

var _ ports.PasswordManager = (*PasswordManager)(nil)

func NewPasswordManager() *PasswordManager {
	return &PasswordManager{}
}
