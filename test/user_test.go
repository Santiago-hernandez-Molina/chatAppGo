package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type LoginTest struct {
	Name   string
	Args   map[string]string
	Status int
}

var LoginTestCase []LoginTest = []LoginTest{
	{
		"User log-in successfully",
		map[string]string{
			"email":    "juan@gmail.com",
			"password": "12345678",
		},
		200,
	},
	{
		"User log-in Failed",
		map[string]string{
			"password": "12345678",
		},
		401,
	},
	{
		"User log-in Failed",
		map[string]string{
			"email":    "juan@gmail.com",
			"password": "1234",
		},
		401,
	},
}

func TestLogin(t *testing.T) {
	for _, user := range LoginTestCase {
		w := MakeRequest("POST", "/login", user.Args, false, LoginUser)
		assert.Equal(t, user.Status, w.Code)
	}
}
