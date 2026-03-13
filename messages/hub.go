package messages

import (
	"encoding/json"
	"log"
)

func (h *Hub) broadcastOnlineUsers() {
	var users []string
	for username := range h.Clients {
		users = append(users, username)
	}

	log.Printf("Broadcasting online users: %v\n", users) // ← ICI

	data, _ := json.Marshal(map[string]interface{}{
		"type":  "online_users",
		"users": users,
	})

	for _, client := range h.Clients {
		log.Printf("Sending to: %s\n", client.UserName) // ← ICI
		client.Send <- data
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
				close(client.Send)
				h.broadcastOnlineUsers()
				log.Printf("User disconnected: %s\n", client.UserName)
			}
		}
	}
}
