package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/config"
	datamongo "github.com/Santiago-hernandez-Molina/chatAppBackend/test/data/data_mongo"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	LoginUser = map[string]string{
		"email":    "juan@gmail.com",
		"password": "12345678",
	}
	LoginUser2 = map[string]string{
		"email":    "pedro@gmail.com",
		"password": "12345678",
	}
	App *gin.Engine
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error reading ENV file")
	}
	config.EMAIL_HOST = os.Getenv("EMAIL_HOST")
	config.EMAIL_USER = os.Getenv("EMAIL_HOST_USER")
	config.EMAIL_PASSWORD = os.Getenv("EMAIL_HOST_PASSWORD")
	config.SECRET = os.Getenv("TEST_SECRET")

	config.MONGO_URI = os.Getenv("TEST_MONGO_URI")
	config.DATABASE_NAME = os.Getenv("TEST_DATABASE_NAME")
	App = config.ConfigApp()
}

func TestMain(m *testing.M) {
	datamongo.InitDB(config.MONGO_URI, config.DATABASE_NAME)
	exitCode := m.Run()

	datamongo.CleanDB()
	os.Exit(exitCode)
}

func MakeRequest[Body map[string]string | map[string]int](
	method,
	url string,
	body Body,
	isAuthenticatedRequest bool,
	user map[string]string,
) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if isAuthenticatedRequest {
		request.AddCookie(authCookie(user))
	}
	writer := httptest.NewRecorder()
	App.ServeHTTP(writer, request)
	return writer
}

func authCookie(userAuth map[string]string) *http.Cookie {
	writer := MakeRequest("POST", "/login", userAuth, false, LoginUser)
	cookies := writer.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "Authorization" {
			return cookie
		}
	}
	return nil
}
