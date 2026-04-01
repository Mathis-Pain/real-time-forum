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
		Send:     make(chan []byte, 256),
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
		// ajout du switch case pour typing
		switch incoming.Type {
		case "private_message":
			handlePrivateMessage(c, incoming)
		case "typing":
			handleTyping(c, incoming)
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
	query := `INSERT INTO messages (SenderID, ReceiverID, Content, CreatedAt) 
              VALUES (?, ?, ?, CURRENT_TIMESTAMP)`
	_, err := db.Exec(query, sender.UserName, in.To, in.Content)
	if err != nil {
		log.Printf("DB insert error: %v\n", err)
		return
	}

	data, err := json.Marshal(OutgoingMessage{
		Type:      "private_message",
		From:      sender.UserName,
		To:        in.To,
		Content:   in.Content,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("JSON marshal error: %v\n", err)
		return
	}

	//  Plus d'accès direct à h.Clients
	HubInstance.Send <- DirectMessage{To: sender.UserName, Data: data}
	HubInstance.Send <- DirectMessage{To: in.To, Data: data}
}

func handleTyping(sender *Client, in IncomingMessage) {
	data, err := json.Marshal(OutgoingMessage{
		Type: "typing",
		From: sender.UserName,
		To:   in.To,
	})
	if err != nil {
		log.Printf("JSON marshal error: %v\n", err)
		return
	}

	//  Uniquement au destinataire
	HubInstance.Send <- DirectMessage{To: in.To, Data: data}
}
