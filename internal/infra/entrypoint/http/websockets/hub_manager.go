package websockets

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type HubManager struct {
	hubs           map[int]*Hub
	roomService    ports.RoomService
	messageService ports.MessageService
	sessionManger  ports.SessionManager
	mutex          *sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHubManager(
	sessionManger ports.SessionManager,
	roomService ports.RoomService,
	messageService ports.MessageService,
) *HubManager {
	return &HubManager{
		hubs:           make(map[int]*Hub),
		mutex:          &sync.Mutex{},
		sessionManger:  sessionManger,
		roomService:    roomService,
		messageService: messageService,
	}
}

func (manager *HubManager) HandleHubs(ctx *gin.Context) {
	cookieAuth, _ := ctx.Request.Cookie("Authorization")
	claims, _ := manager.sessionManger.GetCredentials(cookieAuth.Value)
	var wg sync.WaitGroup

	roomParam := ctx.Param("roomid")
	roomId, _ := strconv.Atoi(roomParam)

	_, err := manager.roomService.GetRoomById(roomId)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusNotFound)
		return
	}

	manager.mutex.Lock()
	wg.Add(2)
	hub := manager.hubs[roomId]
	if hub == nil {
		manager.hubs[roomId] = NewHub(roomId)
		go manager.hubs[roomId].Run(&wg)
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	defer conn.Close()

	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	client := NewClient(
		claims.UserId,
        &models.User{Id: claims.UserId, Email: claims.Email, Username: claims.Username},
		manager.hubs[roomId],
		conn, manager.messageService,
	)

	manager.hubs[roomId].Register <- client
	go client.Run(&wg)
	manager.mutex.Unlock()
	wg.Wait()
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	delete(manager.hubs, roomId)
}
