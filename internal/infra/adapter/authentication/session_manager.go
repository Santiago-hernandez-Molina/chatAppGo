package authentication

import (
	"time"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/golang-jwt/jwt/v5"
)

type SessionManager struct {
	secret string
}

func (manager *SessionManager) SignToken(user *models.User) (string, error) {
	claims := models.UserClaims{
		UserId:   user.Id,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(5 * time.Hour),
			),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSTR, err := token.SignedString([]byte(manager.secret))
	if err != nil {
		return "", err
	}
	return tokenSTR, nil
}

func (manager *SessionManager) GetCredentials(token string) (*models.UserClaims, error) {
	claims := models.UserClaims{}
	_, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(manager.secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &claims, nil
}

var _ ports.SessionManager = (*SessionManager)(nil)

func NewSessionManager(secret string) *SessionManager {
	return &SessionManager{
		secret: secret,
	}
}
