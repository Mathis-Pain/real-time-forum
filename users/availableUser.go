package users

import (
	"database/sql"
	"fmt"
)

// Retourne (true, nil) si disponibles
// Retourne (false, error) si déjà utilisés ou erreur DB
func IsAvailableUser(db *sql.DB, username string, email string) (bool, error) {
	var existingUserName string

	// Vérifier si le username existe
	err := db.QueryRow("SELECT UserName FROM users WHERE UserName = ?", username).Scan(&existingUserName)
	if err == nil {
		// Username trouvé  déjà utilisé
		return false, fmt.Errorf("username already exists")
	}
	if err != sql.ErrNoRows {
		// Vraie erreur de base de données
		return false, fmt.Errorf("database error checking username: %v", err)
	}

	// Vérifier si l'email existe
	var existingEmail string
	err = db.QueryRow("SELECT Email FROM users WHERE Email = ?", email).Scan(&existingEmail)
	if err == nil {
		// Email trouvé déjà utilisé
		return false, fmt.Errorf("email already exists")
	}
	if err != sql.ErrNoRows {
		// Vraie erreur de base de données
		return false, fmt.Errorf("database error checking email: %v", err)
	}

	// Username et email disponibles
	return true, nil
}
