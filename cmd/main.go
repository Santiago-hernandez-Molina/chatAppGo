package main

import (
	"context"
	"log"
	"os"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/services"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/adapter/authentication"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/adapter/repository/mongo"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/handlers"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/middlewares"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/websockets"
	"github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error reading env")
    }
	SECRET := os.Getenv("SECRET")
    USERDB := os.Getenv("USER_DB")
    PASSWORDDB := os.Getenv("PASSWORD_DB")
	ctx := context.Background()

	// repos
    cs := "mongodb://"+USERDB+":"+PASSWORDDB+"@localhost:27017/?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"
	mongoRepo, err := mongo.NewRepo(cs)
	if err != nil {
		log.Fatal(err)
	}
	roomRepo := mongo.NewRoomRepo(mongoRepo, ctx)
	messageRepo := mongo.NewMessageRepo(mongoRepo, ctx)
	userRepo := mongo.NewUserRepo(mongoRepo, ctx)

	// AuthManager
	sessionManager := authentication.NewSessionManager(SECRET)
	passwordManager := authentication.NewPasswordManager()

	// services
	roomService := services.NewRoomService(roomRepo)
	messageService := services.NewMessageService(messageRepo)
	userService := services.NewUserService(
		userRepo,
		sessionManager,
		passwordManager,
	)

	// middlewares
	authMiddleware := middlewares.NewAuthMiddleware(sessionManager)
	roomAccess := middlewares.NewRoomAccess(
		sessionManager,
		roomService,
	)

	// handlers
	roomHandlers := handlers.NewRoomHandler(
		sessionManager,
		roomService,
	)
	userHandlers := handlers.NewUserHandler(
		sessionManager,
		userService,
	)
	messageHandler := handlers.NewMessageHandler(
		messageService,
	)
	hubManager := websockets.NewHubManager(
		sessionManager,
		roomService,
		messageService,
	)

	// httpServer
	router := http.NewServer(
		roomHandlers,
		authMiddleware,
		roomAccess,
		userHandlers,
		messageHandler,
		hubManager,
	)
	router.Start()
}
