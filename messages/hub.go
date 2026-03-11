package messages

import (
	"log"
)

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
