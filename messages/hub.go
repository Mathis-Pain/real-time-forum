package messages

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserName string
	Conn     *websocket.Conn
	Send     chan []byte
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
}

var HubInstance = Hub{
	Clients:    make(map[string]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func (h *Hub) broadcastOnlineUsers() {
	var users []string

	for username := range h.Clients {
		users = append(users, username)
	}

	message := map[string]interface{}{
		"type":  "online_users",
		"users": users,
	}

	for _, client := range h.Clients {
		client.Conn.WriteJSON(message)
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserName] = client
			h.broadcastOnlineUsers()
			log.Printf("User connected: %s\n", client.UserName)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserName]; ok {
				delete(h.Clients, client.UserName)
				h.broadcastOnlineUsers()
				close(client.Send)
				log.Printf("User disconnected: %s\n", client.UserName)
			}
		}
	}
}
