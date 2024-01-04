package websockets

import (
	"log"
	"sync"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
	"github.com/gorilla/websocket"
)

type ClientWS struct {
	Id             int
	User           *models.User
	Conn           *websocket.Conn
	Hub            *Hub
	Outbound       chan []byte
	Inbound        chan RequestWS
	Disconnect     chan int
	MessageService ports.MessageService
}

func NewClient(
	id int,
	user *models.User,
	hub *Hub,
	conn *websocket.Conn,
	messageServ ports.MessageService,
) *ClientWS {
	return &ClientWS{
		Id:             id,
		User:           user,
		Hub:            hub,
		Conn:           conn,
		Outbound:       make(chan []byte),
		Inbound:        make(chan RequestWS),
		Disconnect:     make(chan int),
		MessageService: messageServ,
	}
}

func (client *ClientWS) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	go client.Read()
	for {
		select {
		case message := <-client.Inbound:
			client.Hub.Inbound <- message
		case message := <-client.Outbound:
			client.Conn.WriteMessage(
				websocket.TextMessage,
				message,
			)
		case <-client.Disconnect:
			close(client.Inbound)
			close(client.Outbound)
			close(client.Disconnect)
			client.Conn.Close()
			client.Hub.Unregister <- client
			return
		}
	}
}

func (client *ClientWS) Read() {
	messageUser := models.MessageUser{}
	for {
		err := client.Conn.ReadJSON(&messageUser)
		if err != nil {
			log.Print("error JSON: ")
			log.Println(err)
			client.Disconnect <- 1
			return
		}
		message := models.Message{
			UserId:  client.User.Id,
			RoomId:  client.Hub.Id,
			Content: messageUser.Content,
		}
		messageUser.User = client.User
		err = client.MessageService.SaveMessage(&message)
		if err != nil {
			client.Disconnect <- 1
			log.Println(err)
			return
		}
		client.Hub.Inbound <- RequestWS{
			Client:  client,
			Message: messageUser,
		}
	}
}
