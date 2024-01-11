package config

import (
	"context"
	"log"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/adapter/authentication"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/adapter/email"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/adapter/repository/mongo"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/adapter/tasks"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/handlers"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/middlewares"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/websockets/chats"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/usecases"
	"github.com/gin-gonic/gin"
)

var (
	SECRET         string
	EMAIL_HOST     string
	EMAIL_USER     string
	EMAIL_PASSWORD string
	MONGO_URI      string
	DATABASE_NAME  string
)

func ConfigApp() *gin.Engine {
	// CONTEXT
	ctx := context.Background()

	// REPOSITORIES
	mongoRepo, err := mongo.NewRepo(MONGO_URI, DATABASE_NAME)
	if err != nil {
		log.Fatal("err", err)
	}
	roomRepo := mongo.NewRoomRepo(mongoRepo, ctx)
	messageRepo := mongo.NewMessageRepo(mongoRepo, ctx)
	userRepo := mongo.NewUserRepo(mongoRepo, ctx)
	contactRequestRepo := mongo.NewContactRequestRepo(mongoRepo, ctx)

	// EMAIL SENDER
	emailSender := email.NewEmailSender(
		EMAIL_USER,
		EMAIL_PASSWORD,
		EMAIL_HOST,
	)

	// AUTH MANAGER
	sessionManager := authentication.NewSessionManager(SECRET)
	passwordManager := authentication.NewPasswordManager()

	// TASKS
	userTask := tasks.NewUserTasks(userRepo)

	// USE CASES
	roomUseCase := usecases.NewRoomUseCase(roomRepo, userRepo)
	messageUseCase := usecases.NewMessageUseCase(messageRepo)
	contactRequestUseCase := usecases.NewContactRequestUseCase(
		contactRequestRepo,
		roomRepo,
	)
	userUseCase := usecases.NewUserUseCase(
		userRepo,
		sessionManager,
		passwordManager,
		emailSender,
		userTask,
	)

	// MIDDLEWARES
	authMiddleware := middlewares.NewAuthMiddleware(sessionManager)
	roomAccess := middlewares.NewRoomAccess(
		sessionManager,
		roomUseCase,
	)

	// HANDLERS
	roomManager := chats.NewRoomManager()
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
	contactRequestHandler := handlers.NewContactRequestHandler(
		contactRequestUseCase,
		sessionManager,
	)

	// HTTP ROUTER
	router := http.NewServer(
		roomHandlers,
		authMiddleware,
		roomAccess,
		userHandlers,
		messageHandler,
		contactRequestHandler,
	)
	return router.SetupServer()
}
