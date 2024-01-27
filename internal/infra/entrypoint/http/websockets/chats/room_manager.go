package chats

import (
	"net/http"
	"sync"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type RoomManager struct {
	Hubs     map[int]*Hub
	mutex    *sync.Mutex
	upgrader websocket.Upgrader
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		Hubs:  map[int]*Hub{},
		mutex: &sync.Mutex{},
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (manager *RoomManager) AddClient(
	context *gin.Context,
	user *models.User,
	roomId int,
	hub *Hub,
	messageUseCase ports.MessageUseCase,
) (*Client, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	conn, err := manager.upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		return nil, err
	}
	client := NewClient(conn, hub, user, messageUseCase)
	manager.Hubs[roomId].newClient <- client
	return client, nil
}

func (manager *RoomManager) AddHub(roomId int) *Hub {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if manager.Hubs[roomId] != nil {
		return nil
	}
	manager.Hubs[roomId] = NewHub(roomId)
	return manager.Hubs[roomId]
}

func (manager *RoomManager) RemoveHub(roomId int) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.Hubs, roomId)
}
