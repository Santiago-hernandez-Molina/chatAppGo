package test

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/config"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/test/data"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	App  *gin.Engine
	Auth *http.Cookie

	MONGO_URI     string
	DATABASE_NAME string
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

	MONGO_URI = os.Getenv("TEST_MONGO_URI")
	DATABASE_NAME = os.Getenv("TEST_DATABASE_NAME")

	config.MONGO_URI = MONGO_URI
	config.DATABASE_NAME = DATABASE_NAME
	App = config.ConfigApp()
}

func TestApp(t *testing.T) {
	data.InitDB(MONGO_URI, DATABASE_NAME)
	log.Println("test")
	t.Run("Test Login Endpoint", TLogin)

	defer t.Cleanup(func() {
		data.CleanDB()
	})
}
