package websockets

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
)

type RequestWS struct {
	Client  *ClientWS
	Message models.MessageUser
}


type Hub struct {
	Id         int
	Clients    map[string]*ClientWS
	Register   chan *ClientWS
	Unregister chan *ClientWS
	Inbound    chan RequestWS
	Disconnect chan int
	Mutex      *sync.Mutex
}

func NewHub(roomId int) *Hub {
	return &Hub{
		Id:         roomId,
		Clients:    map[string]*ClientWS{},
		Register:   make(chan *ClientWS),
		Unregister: make(chan *ClientWS),
		Disconnect: make(chan int),
		Mutex:      &sync.Mutex{},
		Inbound:    make(chan RequestWS),
	}
}

func (h *Hub) onConnect(client *ClientWS) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	log.Println(client.Conn.RemoteAddr())

	h.Clients[client.Conn.RemoteAddr().String()] = client
}

func (h *Hub) onDisconnect(client *ClientWS) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	delete(h.Clients, client.Conn.RemoteAddr().String())
}

func (h *Hub) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case client := <-h.Register:
			h.onConnect(client)
		case client := <-h.Unregister:
			h.onDisconnect(client)
			if len(h.Clients) == 0 {
				close(h.Register)
				close(h.Unregister)
				close(h.Inbound)
				return
			}
		case message := <-h.Inbound:
			h.Broadcast(message.Message, message.Client)
		}
	}
}

func (h *Hub) Broadcast(message any, ignore *ClientWS) {
	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	for _, client := range h.Clients {
		if client == ignore {
			continue
		}
		client.Outbound <- data
	}
}
