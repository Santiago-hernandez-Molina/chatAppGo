package utils

import (
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func GetToken(tokenString string, secret string) (*models.UserClaims, error) {
	claims := models.UserClaims{}
	_, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	return &claims, nil
}
