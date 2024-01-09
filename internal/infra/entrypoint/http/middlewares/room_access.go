package middlewares

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/gin-gonic/gin"
)

type RoomAccess struct {
	sessionManager ports.SessionManager
	roomUseCase    ports.RoomUseCase
}

func (middleware *RoomAccess) VerifyRoomAccess(roles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookieAuth, _ := ctx.Request.Cookie("Authorization")
		claims, _ := middleware.sessionManager.GetCredentials(cookieAuth.Value)

		userId := claims.UserId
		roomParam := ctx.Param("roomid")
		roomId, err := strconv.Atoi(roomParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "incorrect room provided",
			})
			ctx.Abort()
			return
		}
		userRoom, err := middleware.roomUseCase.GetUserRoom(
			userId,
			roomId,
		)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "You are not allowed to this room",
			})
			ctx.Abort()
			return
		}
		if !slices.Contains(roles, userRoom.Role) {
			ctx.JSON(http.StatusForbidden, gin.H{
				"message": "You are not allowed to this room",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func NewRoomAccess(sessionManager ports.SessionManager, roomService ports.RoomUseCase) *RoomAccess {
	return &RoomAccess{
		sessionManager: sessionManager,
		roomUseCase:    roomService,
	}
}
