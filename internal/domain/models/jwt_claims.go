package models

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserId   int
	Username string
	Email    string
	jwt.RegisteredClaims
}
