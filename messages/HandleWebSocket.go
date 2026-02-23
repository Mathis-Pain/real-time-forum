package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/auth"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients      = make(map[int]*Client)
	clientsMutex sync.RWMutex
	broadcast    = make(chan IncomingMsg)
)

var database *sql.DB

func HandleWebSocket(db *sql.DB) http.HandlerFunc {
	database = db
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.UserIDKey).(int)
		var name string
		err := db.QueryRow("SELECT UserName FROM users WHERE id=?", userID).Scan(&name)
		if err != nil {
			log.Println("Erreur récupération nom:", err)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Erreur WebSocket upgrade:", err)
			return
		}

		// Mettre en ligne dans la BDD
		_, err = db.Exec("UPDATE users SET userOnline = 1 WHERE id = ?", userID)
		if err != nil {
			log.Printf("Erreur mise à jour statut connexion: %v\n", err)
		} else {
			log.Printf("%s (ID: %d) est maintenant EN LIGNE\n", name, userID)
		}

		// Ajouter le client
		clientsMutex.Lock()
		clients[userID] = &Client{Conn: conn, Name: name}
		totalClients := len(clients)
		clientsMutex.Unlock()

		fmt.Printf("Client connecté: %s (ID: %d) | Total: %d\n", name, userID, totalClients)

		// Délai pour que le client soit prêt avant d'envoyer
		time.Sleep(100 * time.Millisecond)
		BroadcastOnlineUsers()

		// Nettoyage à la déconnexion
		defer func() {
			_, err := db.Exec("UPDATE users SET userOnline = 0 WHERE id = ?", userID)
			if err != nil {
				log.Printf("Erreur mise à jour statut déconnexion: %v\n", err)
			} else {
				log.Printf("%s (ID: %d) est maintenant HORS LIGNE\n", name, userID)
			}

			clientsMutex.Lock()
			delete(clients, userID)
			totalClients := len(clients)
			clientsMutex.Unlock()

			conn.Close()
			fmt.Printf("Client déconnecté: %s (ID: %d) | Total: %d\n", name, userID, totalClients)

			BroadcastOnlineUsers()
		}()

		// Boucle de lecture des messages
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				log.Println("Erreur lecture:", err)
				break
			}

			var incoming IncomingMsg
			if err := json.Unmarshal(msgBytes, &incoming); err != nil {
				log.Println("Erreur JSON:", err)
				continue
			}

			incoming.SenderID = userID
			incoming.SenderName = name

			switch incoming.Type {
			case "message":
				broadcast <- incoming
			case "get_history":
				SendHistory(conn, userID, incoming.ReceiverID, incoming.Offset)
			}
		}
	}
}
