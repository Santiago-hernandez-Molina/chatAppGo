package handlers

import (
	"net/http"
	"strconv"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/infra/entrypoint/http/dtos"
	"github.com/gin-gonic/gin"
)


type RoomHandler struct {
	sessionManager ports.SessionManager
	roomService     ports.RoomService
	messagesService ports.MessageService
}

var _ ports.RoomHandler = (*RoomHandler)(nil)

func NewRoomHandler(
	sessionManager ports.SessionManager,
	roomService ports.RoomService,
) *RoomHandler {
	return &RoomHandler{
		roomService:     roomService,
		sessionManager: sessionManager,
	}
}

func (rh *RoomHandler) GetRoomById(ctx *gin.Context) {
	roomParam := ctx.Param("roomid")
	roomId, _ := strconv.Atoi(roomParam)

	room, err := rh.roomService.GetRoomById(roomId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Room not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, room)
}

func (rh *RoomHandler) AddUserToRoom(ctx *gin.Context) {
	panic("not implemented")
}

func (handler *RoomHandler) GetRoomsByUserId(ctx *gin.Context) {
	cookieAuth, _ := ctx.Request.Cookie("Authorization")
	claims, _ := handler.sessionManager.GetCredentials(cookieAuth.Value)
	userId := claims.UserId

	rooms, err := handler.roomService.GetRoomsByUserId(userId)
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
	err = handler.roomService.NewRoom(
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
