package auth

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type contextKey string

const UserIDKey contextKey = "userID"

// Le middleware intercepte chaque requête vers une route protégée et vérifie la validité d'une session
// Centralisation de la logique d'authentification
// Signature standard d'un middleware en Go
func AuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Lecture du cookie, si il n'existe pas, on laisse passer la requête sans userID (certaines routes publiques peuvent utiliser le même middleware sans être bloquées)
			cookie, err := r.Cookie("session_token")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			token := cookie.Value

			// Vérification de la session en base
			var userID int
			var expiresAt time.Time

			err = db.QueryRow(`
                SELECT UserID, ExpiresAt 
                FROM session 
                WHERE Token = ?
            `, token).Scan(&userID, &expiresAt)

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// Vérifier l'expiration
			if time.Now().After(expiresAt) {
				db.Exec("DELETE FROM session WHERE Token = ?", token)
				next.ServeHTTP(w, r)
				return
			}

			// Marquer l'utilisateur comme online
			db.Exec("UPDATE users SET userOnline = 1 WHERE id = ?", userID)

			// Ajouter au contexte
			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
