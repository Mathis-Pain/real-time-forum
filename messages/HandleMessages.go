package messages

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// HandleMessages écoute en permanence le channel "broadcast".
// chaque message reçu est sauvegardé en base de données,envoyé au destinataire, envoyé a l'expediteur
func HandleMessages() {

	// Boucle infinie : le serveur attend constamment des messages
	for {

		// Attend qu’un message arrive sur le channel "broadcast"
		msg := <-broadcast

		// On récupère l'heure actuelle
		now := time.Now()

		// Insertion du message dans la table messages
		_, err := database.Exec(`
            INSERT INTO messages (SenderID, ReceiverID, Content, CreatedAt)
            VALUES (?, ?, ?, ?)`,
			msg.SenderID,
			msg.ReceiverID,
			msg.Content,
			now,
		)

		// Si erreur SQL, on log et on passe au message suivant
		if err != nil {
			log.Println("Erreur insertion message:", err)
			continue
		}

		fmt.Printf("Message sauvegardé: %s → ID%d : %s\n",
			msg.SenderName,
			msg.ReceiverID,
			msg.Content,
		)

		// Formatage de la date en format standard JSON (ISO 8601)
		createdAt := now.Format(time.RFC3339)

		// Message pour le DESTINATAIRE
		// IsMine = false car ce n’est pas son message
		receiverMsg, _ := json.Marshal(OutgoingMsg{
			Type:      "message",
			Sender:    msg.SenderName,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			CreatedAt: createdAt,
			IsMine:    false,
		})

		// Message pour l’EXPÉDITEUR
		senderMsg, _ := json.Marshal(OutgoingMsg{
			Type:      "message",
			Sender:    msg.SenderName,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			CreatedAt: createdAt,
			IsMine:    true,
		})

		// On verrouille la map en lecture
		clientsMutex.RLock()

		receiver, receiverOnline := clients[msg.ReceiverID]
		sender, senderOnline := clients[msg.SenderID]

		// On libère le verrou
		clientsMutex.RUnlock()

		// On envoie seulement s’il est connecté
		if receiverOnline {

			// Verrouillage de la connexion WebSocket
			receiver.mu.Lock()

			err := receiver.Conn.WriteMessage(
				websocket.TextMessage,
				receiverMsg,
			)

			receiver.mu.Unlock()

			if err != nil {
				log.Printf("Erreur envoi destinataire %d: %v\n",
					msg.ReceiverID,
					err,
				)
			}
		}

		// Permet à l'expéditeur de voir son message apparaître
		if senderOnline {

			sender.mu.Lock()

			err := sender.Conn.WriteMessage(
				websocket.TextMessage,
				senderMsg,
			)

			sender.mu.Unlock()

			if err != nil {
				log.Printf("Erreur envoi expéditeur %d: %v\n",
					msg.SenderID,
					err,
				)
			}
		}
	}
}
