package messages

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// Auth via cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	other := r.URL.Query().Get("with")
	if other == "" {
		http.Error(w, "Missing 'with' param", http.StatusBadRequest)
		return
	}

	limit := 10
	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	rows, err := db.Query(`
        SELECT SenderID, ReceiverID, Content, CreatedAt
        FROM messages
        WHERE (SenderID = ? AND ReceiverID = ?)
           OR (SenderID = ? AND ReceiverID = ?)
        ORDER BY CreatedAt DESC, id DESC
        LIMIT ? OFFSET ?
    `, userName, other, other, userName, limit, offset)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []HistoryMessage
	for rows.Next() {
		var m HistoryMessage
		if err := rows.Scan(&m.From, &m.To, &m.Content, &m.CreatedAt); err != nil {
			continue
		}
		messages = append(messages, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
