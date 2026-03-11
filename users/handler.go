package users

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
)

type UserAPI struct {
    ID            int     `json:"id"`
    UserName      string  `json:"nickname"`
    LastMessageAt *string `json:"last_message_at"`
}

func GetAllUsersHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println("📡 Requête /api/users reçue")

        if r.Method != http.MethodGet {
            http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
            return
        }

        // Récupérer le user connecté via le cookie
        cookie, err := r.Cookie("session_token")
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        var currentUser string
        err = db.QueryRow(`
            SELECT u.UserName FROM users u
            JOIN session s ON s.UserID = u.id
            WHERE s.Token = ? AND s.ExpiresAt > CURRENT_TIMESTAMP
        `, cookie.Value).Scan(&currentUser)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        rows, err := db.Query(`
            SELECT 
                u.id,
                u.UserName,
                MAX(m.CreatedAt) AS last_message_at
            FROM users u
            LEFT JOIN messages m 
                ON (m.SenderID = ? AND m.ReceiverID = u.UserName)
                OR (m.SenderID = u.UserName AND m.ReceiverID = ?)
            WHERE u.UserName != ?
            GROUP BY u.id, u.UserName
            ORDER BY last_message_at DESC, u.UserName ASC
        `, currentUser, currentUser, currentUser)
        if err != nil {
            log.Printf("❌ Erreur SQL: %v\n", err)
            http.Error(w, "Erreur serveur", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var users []UserAPI
        for rows.Next() {
            var u UserAPI
            if err := rows.Scan(&u.ID, &u.UserName, &u.LastMessageAt); err != nil {
                log.Printf("⚠️ Erreur scan: %v\n", err)
                continue
            }
            users = append(users, u)
        }

        if users == nil {
            users = []UserAPI{}
        }

        log.Printf("✅ %d utilisateurs trouvés\n", len(users))
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(users)
    }
}