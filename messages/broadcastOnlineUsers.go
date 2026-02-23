package messages

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func BroadcastOnlineUsers() {
	clientsMutex.RLock()
	users := []OnlineUser{}
	for id, client := range clients {
		users = append(users, OnlineUser{ID: id, Name: client.Name})
	}
	clientsMutex.RUnlock()

	usersJSON, err := json.Marshal(map[string]interface{}{
		"type":  "online_users",
		"users": users,
	})
	if err != nil {
		log.Println("Erreur marshal users:", err)
		return
	}

	fmt.Printf("Diffusion online_users: %d utilisateurs\n", len(users))

	clientsMutex.RLock()
	defer clientsMutex.RUnlock()

	for id, client := range clients {
		client.mu.Lock()
		err := client.Conn.WriteMessage(websocket.TextMessage, usersJSON)
		client.mu.Unlock()

		if err != nil {
			log.Printf("Erreur envoi à client %d: %v\n", id, err)
		}
	}
}
