package chats

import (
	"sync"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/gorilla/websocket"
)

type Client struct {
	id             string
	conn           *websocket.Conn
	outbound       chan models.MessageUser
	hub            *Hub
	user           *models.User
	messageUseCase ports.MessageUseCase
}

func NewClient(conn *websocket.Conn, hub *Hub, user *models.User, messageUseCase ports.MessageUseCase) *Client {
	return &Client{
		id:             conn.RemoteAddr().String(),
		conn:           conn,
		outbound:       make(chan models.MessageUser),
		hub:            hub,
		user:           user,
		messageUseCase: messageUseCase,
	}
}

func (client *Client) Run() {
	var wg sync.WaitGroup
	wg.Add(2)
	go client.Read(&wg)
	go client.Write(&wg)
	wg.Wait()
}

func (client *Client) Write(wg *sync.WaitGroup) {
	defer wg.Done()
	for message := range client.outbound {
		err := client.conn.WriteJSON(message)
		if err != nil {
			close(client.outbound)
		}
	}
	client.hub.disClient <- client
}

func (client *Client) Read(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		message := models.MessageUser{}
		err := client.conn.ReadJSON(&message)
		if err != nil {
			close(client.outbound)
			client.conn.Close()
			return
		}
		if message.Content == "" {
			continue
		}
		id, err := client.messageUseCase.SaveMessage(&models.Message{
			Content: message.Content,
			UserId:  client.user.Id,
			RoomId:  client.hub.id,
		})
		if err != nil {
			close(client.outbound)
			client.conn.Close()
			return
		}
		message.Id = id
		message.User = client.user
		client.hub.inbound <- MessageWithClientId{
			message:  message,
			clientId: client.id,
		}

	}
}
