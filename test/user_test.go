package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TLogin(t *testing.T) {
	for _, user := range LoginTestCase {
		w := httptest.NewRecorder()
		request, _ := json.Marshal(user.Args)
		req, _ := http.NewRequest(
			"POST",
			"/login",
			bytes.NewBuffer(request),
		)
		App.ServeHTTP(w, req)
		assert.Equal(t, user.Status, w.Code)
	}
}

