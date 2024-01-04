package handlers

import (
	"net/http"
	"strconv"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService ports.MessageService
}

func (mh *MessageHandler) GetMessagesByRoomId(ctx *gin.Context) {
	roomParam := ctx.Param("roomid")
	roomid, err := strconv.Atoi(roomParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot read de room id",
		})
		return
	}

	messages, err := mh.messageService.GetMessagesByRoomId(roomid)
	if err != nil {
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "cannot find the room",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, messages)
}

var _ ports.MessageHandler = (*MessageHandler)(nil)

func NewMessageHandler(messageService ports.MessageService) *MessageHandler {
    return &MessageHandler{
        messageService: messageService,
    }
}
