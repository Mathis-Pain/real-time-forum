package messages

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var db *sql.DB

func Init(database *sql.DB) {
	db = database
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading: %v\n", err)
		return
	}

	// Auth via cookie de session
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("No session cookie found.\n")
		conn.Close()
		return
	}

	var userName string
	err = db.QueryRow(`
        SELECT u.UserName
        FROM users u
        JOIN session s ON s.UserID = u.id
        WHERE s.Token = ? AND s.ExpiresAt > CURRENT_TIMESTAMP
    `, cookie.Value).Scan(&userName)
	if err != nil {
		log.Printf("Invalid or expired session: %v\n", err)
		conn.Close()
		return
	}

	client := &Client{
		UserName: userName,
		Conn:     conn,
		Send:     make(chan []byte),
	}

	HubInstance.Register <- client

	go writePump(client)
	readPump(client)
}

func readPump(c *Client) {
	defer func() {
		HubInstance.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v\n", err)
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(msg, &incoming); err != nil {
			log.Printf("Invalid JSON: %v\n", err)
			continue
		}

		if incoming.Type == "private_message" {
			handlePrivateMessage(c, incoming)
		}
	}
}

func writePump(c *Client) {
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("Write error: %v\n", err)
			return
		}
	}
}

func handlePrivateMessage(sender *Client, in IncomingMessage) {
	// 1. Écrire en BDD
	query := `INSERT INTO messages (SenderID, ReceiverID, Content, CreatedAt) 
              VALUES (?, ?, ?, CURRENT_TIMESTAMP)`
	_, err := db.Exec(query, sender.UserName, in.To, in.Content)
	if err != nil {
		log.Printf("DB insert error: %v\n", err)
		return
	}

	// 2. Construire le message de sortie
	now := time.Now()
	out := OutgoingMessage{
		Type:      "private_message",
		From:      sender.UserName,
		To:        in.To,
		Content:   in.Content,
		CreatedAt: now,
	}

	data, err := json.Marshal(out)
	if err != nil {
		log.Printf("JSON marshal error: %v\n", err)
		return
	}

	// 3. Envoyer au sender
	if client, ok := HubInstance.Clients[sender.UserName]; ok {
		client.Send <- data
	}

	// 4. Envoyer au receiver s'il est connecté
	if client, ok := HubInstance.Clients[in.To]; ok {
		client.Send <- data
	}
}
