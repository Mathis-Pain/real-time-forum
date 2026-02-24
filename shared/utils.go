package shared

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Chiffrement des mots de passes.
func Encryption(password []byte) (hashPassword []byte) {
	// Prend "password" pour le chiffrer.
	hashPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln("Error during encryption : %w", err)
	}

	// Test à effacer par la suite.
	fmt.Println(string(hashPassword))
	return hashPassword
}

func GetUserIdByCookie(r *http.Request, db *sql.DB) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, err
	}

	token := cookie.Value

	sqlQuery := `SELECT userid FROM session WHERE token = ?`
	row := db.QueryRow(sqlQuery, token)

	var userId int

	err = row.Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
