package main

import (
	"context"
	"fmt"
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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if os.Getenv("APP_ENV") != "prod" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("ENV ERROR")
		}
	}
	SECRET := os.Getenv("SECRET")
	MONGO_URI := os.Getenv("MONGO_URI")
	ctx := context.Background()

	// repos
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	URI := fmt.Sprintf("%s", MONGO_URI)
	opts := options.Client().ApplyURI(URI).SetServerAPIOptions(serverAPI)

	mongoRepo, err := mongo.NewRepo(opts)
	if err != nil {
		log.Fatal("err", err)
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
    roomManager := websockets.NewRoomManager()
	roomHandlers := handlers.NewRoomHandler(
		sessionManager,
		roomService,
        roomManager,
        messageService,
	)
	userHandlers := handlers.NewUserHandler(
		sessionManager,
		userService,
	)
	messageHandler := handlers.NewMessageHandler(
		messageService,
	)

	// httpServer
	router := http.NewServer(
		roomHandlers,
		authMiddleware,
		roomAccess,
		userHandlers,
		messageHandler,
	)
	router.Start()
}
