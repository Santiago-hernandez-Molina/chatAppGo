package http

import (
	"os"
	"slices"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	roomAccessMiddleware  *middlewares.RoomAccess
	authMiddleware        *middlewares.AuthMiddleware
	roomHandler           ports.RoomHandler
	userHandler           ports.UserHandler
	messageHandler        ports.MessageHandler
	contactRequestHandler ports.ContactRequestHandler
}

func NewServer(
	roomH ports.RoomHandler,
	authMiddleware *middlewares.AuthMiddleware,
	roomAccessMiddleware *middlewares.RoomAccess,
	userHandler ports.UserHandler,
	messageHandler ports.MessageHandler,
	contactRequestHandler ports.ContactRequestHandler,
) *Server {
	return &Server{
		authMiddleware:        authMiddleware,
		roomHandler:           roomH,
		userHandler:           userHandler,
		messageHandler:        messageHandler,
		roomAccessMiddleware:  roomAccessMiddleware,
		contactRequestHandler: contactRequestHandler,
	}
}

func (server *Server) SetupServer() *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))
	app := gin.Default()
	server.globalMiddlewares(app)

	server.authRoutes(app)

	// API V1
	v1 := app.Group("/v1")
	server.roomRoutes(v1)
	server.messagesRoutes(v1)
	server.contactRequestRoutes(v1)
	server.userRoutes(v1)
	
	return app
}

func (server *Server) globalMiddlewares(app *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://192.168.0.4:5173",
		"https://chatapp-go-vue.netlify.app",
	}
	config.AllowMethods = []string{"POST", "GET"}
	config.AllowCredentials = true
	config.AllowOriginFunc = func(origin string) bool {
		return slices.Contains[[]string](
			[]string{
				"http://localhost:5173",
				"http://192.168.0.4:5173",
				"https://chatapp-go-vue.netlify.app",
			}, origin,
		)
	}

	app.SetTrustedProxies([]string{"localhost"})
	app.Use(cors.New(config))
	app.Use(server.authMiddleware.CheckAuthMiddleware)
}

// AUTH ROUTES
func (server *Server) authRoutes(app *gin.Engine) {
	app.POST("register", server.userHandler.Register)
	app.POST("login", server.userHandler.Login)
	app.POST("activate", server.userHandler.ActivateAccount)
}

// User Routes
func (server *Server) userRoutes(app *gin.RouterGroup) {
	userRoutes := app.Group("/user/")
	userRoutes.GET("/find", server.userHandler.GetUsers)
}

// ContactRequestHandler Routes

func (server *Server) contactRequestRoutes(app *gin.RouterGroup) {
	contactRoutes := app.Group("/contact/")
	contactRoutes.POST("/new", server.contactRequestHandler.SendRequest)
	contactRoutes.GET("/find/received", server.contactRequestHandler.GetReceivedRequests)
	contactRoutes.GET("/find/sended", server.contactRequestHandler.GetSendedRequests)
	contactRoutes.POST("/accept/:requestid", server.contactRequestHandler.AcceptRequest)
}

// Room Routes
func (server *Server) roomRoutes(app *gin.RouterGroup) {
	roomRoutes := app.Group("/room")
	accessUserRoom := roomRoutes.Group(
		"",
		server.roomAccessMiddleware.VerifyRoomAccess([]string{
			"admin",
			"user",
		}),
	)
	var adminRoutes *gin.RouterGroup = roomRoutes.Group(
		"",
		server.roomAccessMiddleware.VerifyRoomAccess([]string{
			"admin",
		}),
	)
	roomRoutes.GET("/find", server.roomHandler.GetRoomsByUserId)
	roomRoutes.POST("/new", server.roomHandler.NewRoom)
	accessUserRoom.GET("/ws/:roomid", server.roomHandler.ConnectToRoom)
	accessUserRoom.GET("/show/:roomid", server.roomHandler.GetRoomById)
	adminRoutes.POST("/:roomid/add-user", server.roomHandler.AddUserToRoom)
}

// Messages Routes
func (server *Server) messagesRoutes(app *gin.RouterGroup) {
	messageRoutes := app.Group("/message")
	accessRoomMessages := messageRoutes.Group(
		"",
		server.roomAccessMiddleware.VerifyRoomAccess([]string{
			"admin",
			"user",
		}),
	)
	accessRoomMessages.GET("/find/room/:roomid", server.messageHandler.GetMessagesByRoomId)
}
