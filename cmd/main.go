package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"real-time-forum/auth"
	"real-time-forum/internal/config-database" // Ceci permet d'appeler la fonction qui se trouve dans le fichier.
	"real-time-forum/messages"
	"real-time-forum/server"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// MARK: .ENV
	// Chargement des données du .env
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
		return
	}
	// Attribution du port du serveur.
	// port := os.Getenv("SERVER_PORT")
	port := ":8080" // Inscrit en dur pour les tests.
	// Attribution du chemin de la database.
	// pathDB := os.Getenv("REALTIMEFORUM_DB_PATH")
	pathDB := "./vault/real_time_forum_database.db" // Inscrit en dur pour les tests.

	// MARK: DB
	// Lancement de la BDD.
	// Vérification si la BDD existe déjà.
	// Créer le dossier vault s'il n'existe pas
	if err := os.MkdirAll("./vault", 0755); err != nil {
		log.Fatalf("Error creating vault directory: %v", err)
	}

	var db *sql.DB
	var err error

	_, statErr := os.Stat(pathDB) // Stat repère si le fichier existe, sans le charger.
	if statErr != nil {
		// Création de la BDD.
		fmt.Println("Initialazing Database...")
		db, err = config.RunDB(pathDB)
		db, err = config.InitDB(pathDB, db)
		config.InspectDbIntegrity(db)
		fmt.Println("Connection to Database...")
	} else {
		// Ouverture de la BDD.
		fmt.Println("Connection to Database...")
		db, err = config.RunDB(pathDB)
		config.InspectDbIntegrity(db)
	}
	messages.Init(db)
	go messages.HubInstance.Run()
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer db.Close()
	// MARK: Server
	// Nettoyage immédiat des sessions au démarrage
	if err := auth.CleanExpiredSessions(db); err != nil {
		log.Println("Erreur nettoyage initial des sessions :", err)
	} else {
		fmt.Println("Nettoyage initial des sessions terminé")
	}

	// Nettoyage automatique toutes les 5 minutes
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // toutes les 5 minutes
		for {
			<-ticker.C
			err := auth.CleanExpiredSessions(db)
			if err != nil {
				log.Println("Erreur nettoyage sessions :", err)
			} else {
				fmt.Println("Nettoyage des sessions expirées terminé")
			}
		}
	}()
	// MARK: WebSocket

	fmt.Println("WebSocket disponible sur ws://localhost" + port + "/ws")
	// Routes HTTP normales (à adapter selon votre code)
	// http.HandleFunc("/", yourHandler)
	// http.HandleFunc("/login", yourLoginHandler)
	// etc...

	server.Server(port, db)
}
