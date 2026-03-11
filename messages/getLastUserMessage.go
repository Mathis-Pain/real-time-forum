package messages

import (
	"encoding/json"
	"net/http"
)

func GetLastUserMessage(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var userID int
	var toUserName string

	err = db.QueryRow(`
        SELECT u.id, u.UserName
        FROM users u
        JOIN session s ON s.UserID = u.id
        WHERE s.Token = ? AND s.ExpiresAt > CURRENT_TIMESTAMP
    `, cookie.Value).Scan(&userID, &toUserName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query(`
        SELECT u.UserName, CreatedAt
        FROM messages m
				INNER JOIN users u ON u.id = m.SenderID
        WHERE m.ReceiverID = ? OR m.SenderID = ?
				ORDER BY m.CreatedAt DESC
    `, userID, userID)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var listLastUserMessage []lastUserMessage

	for rows.Next() {
		var lastUser lastUserMessage
		var senderUserName string

		if err = rows.Scan(&senderUserName, &lastUser.CreatedAt); err != nil {
			continue
		}

		lastUser.From = senderUserName
		lastUser.To = toUserName

		listLastUserMessage = append(listLastUserMessage, lastUser)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listLastUserMessage)
}
