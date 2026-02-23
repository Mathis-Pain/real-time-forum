package messages

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func SendHistory(conn *websocket.Conn, userID int, otherID int, offset int) {
	// Récupérer 11 messages pour détecter s'il y en a plus
	rows, err := database.Query(`
        SELECT SenderID, Content, CreatedAt
        FROM messages
        WHERE (SenderID = ? AND ReceiverID = ?)
           OR (SenderID = ? AND ReceiverID = ?)
    ORDER BY datetime(CreatedAt) DESC
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

		// Chercher le nom d'abord dans les clients connectés
		senderName := ""
		clientsMutex.RLock()
		if client, ok := clients[senderID]; ok {
			senderName = client.Name
		}
		clientsMutex.RUnlock()

		// Sinon chercher en BDD
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

	// Détecter s'il y a plus de messages
	hasMore := len(messages) == 11
	if hasMore {
		messages = messages[:10] // Garder seulement 10
	}

	// Inverser : les plus récents en bas
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
