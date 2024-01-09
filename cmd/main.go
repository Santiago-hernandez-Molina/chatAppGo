package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/adapter/authentication"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/adapter/email"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/adapter/repository/mongo"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/adapter/tasks"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/handlers"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/middlewares"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/websockets"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/usecases"
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

	EMAIL_HOST := os.Getenv("EMAIL_HOST")
	EMAIL_USER := os.Getenv("EMAIL_HOST_USER")
	EMAIL_PASSWORD := os.Getenv("EMAIL_HOST_PASSWORD")

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

	// Email Sender
	emailSender := email.NewEmailSender(
		EMAIL_USER,
		EMAIL_PASSWORD,
		EMAIL_HOST,
	)

	// AuthManager
	sessionManager := authentication.NewSessionManager(SECRET)
	passwordManager := authentication.NewPasswordManager()

	// Tasks
	userTask := tasks.NewUserTasks(userRepo)

	// UseCases
	roomUseCase := usecases.NewRoomUseCase(roomRepo)
	messageUseCase := usecases.NewMessageUseCase(messageRepo)
	userUseCase := usecases.NewUserUseCase(
		userRepo,
		sessionManager,
		passwordManager,
		emailSender,
		userTask,
	)

	// Middlewares
	authMiddleware := middlewares.NewAuthMiddleware(sessionManager)
	roomAccess := middlewares.NewRoomAccess(
		sessionManager,
		roomUseCase,
	)

	// Handlers
	roomManager := websockets.NewRoomManager()
	roomHandlers := handlers.NewRoomHandler(
		sessionManager,
		roomUseCase,
		roomManager,
		messageUseCase,
	)
	userHandlers := handlers.NewUserHandler(
		sessionManager,
		userUseCase,
		userTask,
	)
	messageHandler := handlers.NewMessageHandler(
		messageUseCase,
	)

	// HttpServer
	router := http.NewServer(
		roomHandlers,
		authMiddleware,
		roomAccess,
		userHandlers,
		messageHandler,
	)
	router.Start()
}
