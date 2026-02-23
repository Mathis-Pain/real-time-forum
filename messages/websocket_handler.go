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

type Client struct {
	Conn *websocket.Conn
	Name string
	mu   sync.Mutex
}

var (
	clients      = make(map[int]*Client)
	clientsMutex sync.RWMutex
	broadcast    = make(chan IncomingMsg)
)

type IncomingMsg struct {
	Type       string `json:"type"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	Offset     int    `json:"offset"`
	SenderID   int
	SenderName string
}

// IsMine ajouté pour l'expéditeur ET le destinataire
type OutgoingMsg struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	IsMine    bool   `json:"is_mine"`
}

type HistoryMsg struct {
	Sender    string `json:"sender"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	IsMine    bool   `json:"is_mine"`
}

type OnlineUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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
		broadcastOnlineUsers()

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

			broadcastOnlineUsers()
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
				sendHistory(conn, userID, incoming.ReceiverID, incoming.Offset)
			}
		}
	}
}

func broadcastOnlineUsers() {
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

func HandleMessages() {
	for {
		msg := <-broadcast
		now := time.Now()

		// Sauvegarder en BDD
		_, err := database.Exec(`
            INSERT INTO messages (SenderID, ReceiverID, Content, CreatedAt)
            VALUES (?, ?, ?, ?)`,
			msg.SenderID, msg.ReceiverID, msg.Content, now,
		)
		if err != nil {
			log.Println("Erreur insertion message:", err)
			continue
		}

		fmt.Printf("Message sauvegardé: %s → ID%d : %s\n", msg.SenderName, msg.ReceiverID, msg.Content)

		createdAt := now.Format(time.RFC3339)

		// Envoyer au DESTINATAIRE avec is_mine = false
		receiverMsg, _ := json.Marshal(OutgoingMsg{
			Type:      "message",
			Sender:    msg.SenderName,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			CreatedAt: createdAt,
			IsMine:    false, // Le destinataire reçoit le message de quelqu'un d'autre
		})

		// Envoyer à l'EXPÉDITEUR avec is_mine = true
		senderMsg, _ := json.Marshal(OutgoingMsg{
			Type:      "message",
			Sender:    msg.SenderName,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			CreatedAt: createdAt,
			IsMine:    true, // L'expéditeur voit son propre message
		})

		clientsMutex.RLock()
		receiver, receiverOnline := clients[msg.ReceiverID]
		sender, senderOnline := clients[msg.SenderID]
		clientsMutex.RUnlock()

		// Envoyer au destinataire s'il est connecté
		if receiverOnline {
			receiver.mu.Lock()
			err := receiver.Conn.WriteMessage(websocket.TextMessage, receiverMsg)
			receiver.mu.Unlock()
			if err != nil {
				log.Printf("Erreur envoi destinataire %d: %v\n", msg.ReceiverID, err)
			}
		}

		// ✅ Renvoyer à l'expéditeur pour confirmer l'envoi
		if senderOnline {
			sender.mu.Lock()
			err := sender.Conn.WriteMessage(websocket.TextMessage, senderMsg)
			sender.mu.Unlock()
			if err != nil {
				log.Printf("Erreur envoi expéditeur %d: %v\n", msg.SenderID, err)
			}
		}
	}
}

func sendHistory(conn *websocket.Conn, userID int, otherID int, offset int) {
	// ✅ Récupérer 11 messages pour détecter s'il y en a plus
	rows, err := database.Query(`
        SELECT SenderID, Content, CreatedAt
        FROM messages
        WHERE (SenderID = ? AND ReceiverID = ?)
           OR (SenderID = ? AND ReceiverID = ?)
        ORDER BY CreatedAt DESC
        LIMIT 11 OFFSET ?`,
		userID, otherID, otherID, userID, offset,
	)
	if err != nil {
		log.Println("Erreur historique:", err)
		return
	}
	defer rows.Close()

	var messages []HistoryMsg
	for rows.Next() {
		var senderID int
		var content, createdAt string
		rows.Scan(&senderID, &content, &createdAt)

		// ✅ Chercher le nom d'abord dans les clients connectés
		senderName := ""
		clientsMutex.RLock()
		if client, ok := clients[senderID]; ok {
			senderName = client.Name
		}
		clientsMutex.RUnlock()

		// ✅ Sinon chercher en BDD
		if senderName == "" {
			database.QueryRow("SELECT UserName FROM users WHERE id=?", senderID).Scan(&senderName)
		}

		messages = append(messages, HistoryMsg{
			Sender:    senderName,
			SenderID:  senderID,
			Content:   content,
			CreatedAt: createdAt,
			IsMine:    senderID == userID,
		})
	}

	// ✅ Détecter s'il y a plus de messages
	hasMore := len(messages) == 11
	if hasMore {
		messages = messages[:10] // Garder seulement 10
	}

	// ✅ Inverser : les plus récents en bas
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	response, _ := json.Marshal(map[string]interface{}{
		"type":     "message_history",
		"messages": messages,
		"offset":   offset,
		"has_more": hasMore,
	})

	conn.WriteMessage(websocket.TextMessage, response)
}
