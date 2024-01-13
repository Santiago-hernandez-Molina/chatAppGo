package main

import (
	"log"
	"os"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/config"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") != "prod" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error reading env file")
		}
	}

	config.SECRET = os.Getenv("SECRET")
	config.EMAIL_HOST = os.Getenv("EMAIL_HOST")
	config.EMAIL_USER = os.Getenv("EMAIL_HOST_USER")
	config.EMAIL_PASSWORD = os.Getenv("EMAIL_HOST_PASSWORD")
	config.MONGO_URI = os.Getenv("MONGO_URI")
	config.DATABASE_NAME = os.Getenv("DATABASE_NAME")

	app := config.ConfigApp()
	app.Run(":8080")
}
