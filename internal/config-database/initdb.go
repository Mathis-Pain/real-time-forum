package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// MARK: DB opening
// Lancement de la BDD déjà existante.
func RunDB(pathDB string) (*sql.DB, error) {
	// Ouverture ou création de la BDD
	db, err := sql.Open("sqlite3", pathDB)
	if err != nil {
		log.Fatalln("Error opening database : %w", err)
	}
	// Vérification de la connexion, car sql.Open ne le fais pas.
	if err = db.Ping(); err != nil {
		log.Fatalln("Error connecting to database : %w", err)
	}
	return db, nil
}

// MARK: DB Creation
// Création de la BDD si elle n'existe pas.
func InitDB(pathDB string, db *sql.DB) (*sql.DB, error) {
	// Lecture des fichiers contenant les tables.
	// Ouverture du fichier.
	testTables, err := os.ReadFile("./internal/config-database/001_create_tables.sql")
	if err != nil {
		log.Fatalln("Error with testTables in internal/config-database/001_create_tables.sql : %w", err)
	}
	// Ecriture dans le BDD.
	_, err = db.Exec(string(testTables))
	if err != nil {
		log.Fatalln("Error writing testTables in DB : %w", err)
	}
	return db, nil
}
