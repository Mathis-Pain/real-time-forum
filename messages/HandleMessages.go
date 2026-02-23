package messages

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

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

		//  Renvoyer à l'expéditeur pour confirmer l'envoi
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
