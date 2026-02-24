package auth

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func CreateSession(db *sql.DB, userID int) (string, error) {
	token := uuid.NewString()
	now := time.Now()
	expireAt := now.Add(time.Hour * 24)

	_, err := db.Exec("INSERT INTO session (UserID, Token, CreatedAt, ExpiresAt) VALUES (?, ?, ?, ?);", userID, token, now, expireAt)
	if err != nil {
		return "", err
	}

	return token, nil
}

func CleanExpiredSessions(db *sql.DB) error {
	// Supprime toutes les sessions expirées
	_, err := db.Exec("DELETE FROM session WHERE ExpiresAt < CURRENT_TIMESTAMP")
	if err != nil {
		return err
	}

	// Met les utilisateurs sans session active à offline
	_, err = db.Exec("UPDATE users SET userOnline = 0 WHERE id NOT IN (SELECT UserID FROM session)")
	if err != nil {
		return err
	}

	return nil
}

func GetActiveSessionToken(db *sql.DB, userID int) string {
	var sessionToken string
	err := db.QueryRow("SELECT token FROM session WHERE userID = ? AND ExpiresAt > CURRENT_TIMESTAMP", userID).Scan(&sessionToken)
	if err != nil {
		return ""
	}
	return sessionToken
}
