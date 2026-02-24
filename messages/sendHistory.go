package messages

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// SendHistory envoie l’historique des messages entre userID et otherID.
// Paramètres :
//   - conn     : connexion WebSocket du demandeur
//   - userID   : ID de l’utilisateur connecté
//   - otherID  : ID de l’autre utilisateur (conversation)
//   - offset   : nombre de messages déjà chargés (pagination)
//
// La fonction :
//  1. Récupère 11 messages (pour détecter s’il y en a plus)
//  2. Construit la liste des messages formatée pour le frontend
//  3. Indique s’il reste d’autres messages à charger (has_more)
//  4. Envoie le tout en JSON via WebSocket
func SendHistory(conn *websocket.Conn, userID int, otherID int, offset int) {

	//  Requête SQL pour récupérer les messages
	// On récupère 11 messages :
	// - 10 pour affichage
	// - 1 supplémentaire pour savoir s’il y en a encore après
	rows, err := database.Query(`
        SELECT SenderID, Content, CreatedAt
        FROM messages
        WHERE (SenderID = ? AND ReceiverID = ?)
           OR (SenderID = ? AND ReceiverID = ?)
        ORDER BY datetime(CreatedAt) DESC
        LIMIT 11 OFFSET ?`,
		userID, otherID,
		otherID, userID,
		offset,
	)

	if err != nil {
		log.Println("Erreur historique:", err)
		return
	}
	defer rows.Close()

	// 2. Construction de la liste des messages
	var messages []HistoryMsg

	for rows.Next() {

		var senderID int
		var content, createdAt string

		rows.Scan(&senderID, &content, &createdAt)

		// Récupération du nom de l'expéditeur
		senderName := ""

		// D'abord chercher parmi les clients connectés
		clientsMutex.RLock()
		if client, ok := clients[senderID]; ok {
			senderName = client.Name
		}
		clientsMutex.RUnlock()

		// Si non trouvé (utilisateur hors ligne),
		// on cherche dans la base de données
		if senderName == "" {
			database.QueryRow(
				"SELECT UserName FROM users WHERE id=?",
				senderID,
			).Scan(&senderName)
		}

		// Ajout du message formaté dans le tableau
		messages = append(messages, HistoryMsg{
			Sender:    senderName,
			SenderID:  senderID,
			Content:   content,
			CreatedAt: createdAt,

			// Permet au frontend de savoir
			// si le message appartient à l’utilisateur connecté
			IsMine: senderID == userID,
		})
	}

	//  Détection s’il reste des messages
	// Si on a 11 messages, cela signifie
	// qu'il y en a au moins un de plus
	hasMore := len(messages) == 11

	// On garde seulement les 10 premiers pour affichage
	if hasMore {
		messages = messages[:10]
	}

	// Inverser l’ordre des messages
	// La requête SQL renvoie les plus récents en premier.
	// On inverse pour que les plus anciens soient en haut
	// et les plus récents en bas (affichage naturel du chat).
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	// 5. Préparation de la réponse JSON
	response, _ := json.Marshal(map[string]interface{}{
		"type":     "message_history",
		"messages": messages,
		"offset":   offset,
		"has_more": hasMore,
	})

	// Envoi au client via WebSocket
	conn.WriteMessage(websocket.TextMessage, response)
}
