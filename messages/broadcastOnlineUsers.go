package messages

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// BroadcastOnlineUsers envoie à tous les clients connectés
// la liste actuelle des utilisateurs en ligne via WebSocket.
// Le message est envoyé au format JSON avec le type "online_users"
func BroadcastOnlineUsers() {
	//  On verrouille en lecture la map "clients"
	// RLock permet plusieurs lectures simultanées
	clientsMutex.RLock()

	// On crée un tableau vide qui contiendra les utilisateurs connectés
	users := []OnlineUser{}

	// On parcourt la map des clients connectés
	for id, client := range clients {
		// Pour chaque client, on ajoute son ID et son Name dans le tableau
		users = append(users, OnlineUser{
			ID:   id,
			Name: client.Name,
		})
	}

	// On libère le verrou de lecture
	clientsMutex.RUnlock()

	// On transforme les données en JSON
	// Le message envoyé aura cette forme :
	// {
	//   "type": "online_users",
	//   "users": [...]
	// }
	usersJSON, err := json.Marshal(map[string]interface{}{
		"type":  "online_users",
		"users": users,
	})

	// Si erreur pendant la conversion en JSON
	if err != nil {
		log.Println("Erreur marshal users:", err)
		return
	}

	// Log côté serveur pour voir combien d'utilisateurs sont envoyés
	fmt.Printf("Diffusion online_users: %d utilisateurs\n", len(users))

	// On reprend un verrou de lecture pour parcourir les clients
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()

	// On envoie le message à chaque client connecté
	for id, client := range clients {

		// On verrouille le mutex interne du client
		// Important pour éviter que deux goroutines écrivent
		// en même temps sur la même connexion WebSocket
		client.mu.Lock()

		// Envoi du message WebSocket (type texte)
		err := client.Conn.WriteMessage(websocket.TextMessage, usersJSON)

		// On déverrouille
		client.mu.Unlock()

		// Si erreur pendant l’envoi
		if err != nil {
			log.Printf("Erreur envoi à client %d: %v\n", id, err)
		}
	}
}
