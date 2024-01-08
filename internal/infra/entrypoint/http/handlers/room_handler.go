package handlers

import (
	"net/http"
	"strconv"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/dtos"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/websockets"
	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	sessionManager  ports.SessionManager
	roomUseCase     ports.RoomUseCase
	messagesService ports.MessageUseCase
	roomManager     *websockets.RoomManager
}

var _ ports.RoomHandler = (*RoomHandler)(nil)

func NewRoomHandler(
	sessionManager ports.SessionManager,
	roomService ports.RoomUseCase,
	roomManager *websockets.RoomManager,
	messagesService ports.MessageUseCase,
) *RoomHandler {
	return &RoomHandler{
		roomUseCase:     roomService,
		sessionManager:  sessionManager,
		roomManager:     roomManager,
		messagesService: messagesService,
	}
}

func (handler *RoomHandler) ConnectToRoom(ctx *gin.Context) {
	cookieAuth, _ := ctx.Cookie("Authorization")
	roomParam := ctx.Param("roomid")

	claims, _ := handler.sessionManager.GetCredentials(cookieAuth)
	roomId, _ := strconv.Atoi(roomParam)

	hub := handler.roomManager.AddHub(roomId)
	go func() {
		if hub != nil {
			hub.Run()
			handler.roomManager.RemoveHub(roomId)
		}
	}()
	user := &models.User{
		Id:       claims.UserId,
		Username: claims.Username,
		Email:    claims.Email,
	}
	client, err := handler.roomManager.AddClient(
		ctx,
		user,
		roomId,
		handler.roomManager.Hubs[roomId],
		handler.messagesService,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "error connecting ws",
		})
		return
	}
	go client.Run()
}

func (handler *RoomHandler) GetRoomById(ctx *gin.Context) {
	roomParam := ctx.Param("roomid")
	roomId, _ := strconv.Atoi(roomParam)

	room, err := handler.roomUseCase.GetRoomById(roomId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Room not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, room)
}

func (handler *RoomHandler) AddUserToRoom(ctx *gin.Context) {
	userRoom := models.UserRoom{}
	roomParam := ctx.Param("roomid")
	roomId, _ := strconv.Atoi(roomParam)

	err := ctx.BindJSON(&userRoom)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot read the user",
		})
		return
	}
	err = handler.roomUseCase.AddUserToRoom(userRoom.UserId, roomId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error adding new user check your data",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User added sucessfully",
	})
}

func (handler *RoomHandler) GetRoomsByUserId(ctx *gin.Context) {
	cookieAuth, _ := ctx.Request.Cookie("Authorization")
	claims, _ := handler.sessionManager.GetCredentials(cookieAuth.Value)
	userId := claims.UserId

	rooms, err := handler.roomUseCase.GetRoomsByUserId(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Cannot found the rooms",
		})
		return
	}
	ctx.JSON(http.StatusOK, rooms)
}

func (handler *RoomHandler) NewRoom(ctx *gin.Context) {
	cookieAuth, _ := ctx.Request.Cookie("Authorization")
	claims, _ := handler.sessionManager.GetCredentials(cookieAuth.Value)
	userId := claims.UserId

	request := dtos.NewRoomRequest{}
	err := ctx.BindJSON(&request)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = handler.roomUseCase.NewRoom(
		models.Room{Name: request.Name},
		userId,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Cannot Create the room",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Room created sucessfully",
	})
}
