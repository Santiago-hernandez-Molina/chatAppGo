package chats

import (
	"sync"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
)

type MessageWithClientId struct {
	clientId string
	message  models.MessageUser
}

type Hub struct {
	id          int
	clients     map[string]*Client
	inbound     chan MessageWithClientId
	newClient   chan *Client
	disClient   chan *Client
	mutex       *sync.Mutex
	disconnect  chan int
	roomManager *RoomManager
}

func NewHub(roomId int) *Hub {
	return &Hub{
		id:         roomId,
		clients:    make(map[string]*Client),
		inbound:    make(chan MessageWithClientId),
		newClient:  make(chan *Client),
		disClient:  make(chan *Client),
		mutex:      &sync.Mutex{},
		disconnect: make(chan int),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case message := <-hub.inbound:
			hub.Broadcast(message.message, message.clientId)
		case client := <-hub.newClient:
			hub.AddClient(client)
		case client := <-hub.disClient:
			disconnect := hub.RemoveClient(client)
			if disconnect {
				return
			}
		}
	}
}

func (hub *Hub) AddClient(client *Client) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.clients[client.id] = client
}

func (hub *Hub) RemoveClient(client *Client) bool {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	delete(hub.clients, client.id)
	if len(hub.clients) == 0 {
		close(hub.inbound)
		close(hub.newClient)
		close(hub.disClient)
		return true
	}
	return false
}

func (hub *Hub) Broadcast(message models.MessageUser, clientId string) {
	for _, client := range hub.clients {
		if client.id == clientId {
			continue
		}
		client.outbound <- message
	}
}
