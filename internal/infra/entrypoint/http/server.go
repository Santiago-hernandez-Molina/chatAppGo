package http

import (
	"slices"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	authMiddleware       *middlewares.AuthMiddleware
	roomHandler          ports.RoomHandler
	userHandler          ports.UserHandler
	messageHandler       ports.MessageHandler
	roomAccessMiddleware *middlewares.RoomAccess
}

func NewServer(
	roomH ports.RoomHandler,
	authMiddleware *middlewares.AuthMiddleware,
	roomAccessMiddleware *middlewares.RoomAccess,
	userHandler ports.UserHandler,
	messageHandler ports.MessageHandler,
) *Server {
	return &Server{
		authMiddleware:       authMiddleware,
		roomHandler:          roomH,
		userHandler:          userHandler,
		messageHandler:       messageHandler,
		roomAccessMiddleware: roomAccessMiddleware,
	}
}

func (server *Server) Start() {
	app := gin.Default()
	server.globalMiddlewares(app)

	server.authRoutes(app)

	// API V1
	v1 := app.Group("/v1")
	server.roomRoutes(v1)
	server.messagesRoutes(v1)

	app.Run(":8080")
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
	app.POST("/register", server.userHandler.Register)
	app.POST("login", server.userHandler.Login)
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
