package messages

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/auth"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// upgrader permet de transformer une connexion HTTP classique en connexion WebSocket.
// CheckOrigin retourne true pour autoriser toutes les origines.
// En production, il est conseillé de restreindre cela pour la sécurité.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Variables globales du système de messagerie.

var (
	// Map des clients connectés : clé = userID, valeur = pointeur vers Client
	clients = make(map[int]*Client)

	// Mutex pour protéger l'accès concurrent à la map clients
	clientsMutex sync.RWMutex

	// Channel utilisé pour transmettre les messages entrants
	// vers la fonction HandleMessages
	broadcast = make(chan IncomingMsg)
)

// Pointeur global vers la base de données
var database *sql.DB

// HandleWebSocket gère la connexion WebSocket d’un utilisateur:
//   - récupère son identité
//   - ouvre la connexion WebSocket
//   - le marque en ligne en base
//   - l’ajoute à la liste des clients connectés
//   - écoute ses messages en continu
//   - nettoie proprement lors de la déconnexion
func HandleWebSocket(db *sql.DB) http.HandlerFunc {

	// Stocke la connexion base de données dans la variable globale
	database = db

	return func(w http.ResponseWriter, r *http.Request) {

		// Récupération de l’ID utilisateur depuis le contexte (middleware auth)
		userID := r.Context().Value(auth.UserIDKey).(int)

		// Récupération du nom de l’utilisateur en base
		var name string
		err := db.QueryRow(
			"SELECT UserName FROM users WHERE id=?",
			userID,
		).Scan(&name)

		if err != nil {
			log.Println("Erreur récupération nom:", err)
			return
		}

		// Transformation de la requête HTTP en connexion WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Erreur WebSocket upgrade:", err)
			return
		}

		_, err = db.Exec(
			"UPDATE users SET userOnline = 1 WHERE id = ?",
			userID,
		)

		if err != nil {
			log.Printf("Erreur mise à jour statut connexion: %v\n", err)
		} else {
			log.Printf("%s (ID: %d) est maintenant EN LIGNE\n", name, userID)
		}

		// Ajout dans la liste des clients

		clientsMutex.Lock()
		clients[userID] = &Client{
			Conn: conn,
			Name: name,
		}
		totalClients := len(clients)
		clientsMutex.Unlock()

		fmt.Printf(
			"Client connecté: %s (ID: %d) | Total: %d\n",
			name,
			userID,
			totalClients,
		)

		// Petit délai pour laisser le frontend prêt
		time.Sleep(100 * time.Millisecond)

		// Diffuse la nouvelle liste des utilisateurs en ligne
		BroadcastOnlineUsers()

		// Nettoyage à la déconnexion
		defer func() {

			// Mettre l'utilisateur hors ligne en base
			_, err := db.Exec(
				"UPDATE users SET userOnline = 0 WHERE id = ?",
				userID,
			)

			if err != nil {
				log.Printf("Erreur mise à jour statut déconnexion: %v\n", err)
			} else {
				log.Printf("%s (ID: %d) est maintenant HORS LIGNE\n", name, userID)
			}

			// Supprimer de la map des clients connectés
			clientsMutex.Lock()
			delete(clients, userID)
			totalClients := len(clients)
			clientsMutex.Unlock()

			// Fermer la connexion WebSocket
			conn.Close()

			fmt.Printf(
				"Client déconnecté: %s (ID: %d) | Total: %d\n",
				name,
				userID,
				totalClients,
			)

			// Mettre à jour la liste des utilisateurs en ligne
			BroadcastOnlineUsers()
		}()

		// Boucle de lecture des messages
		for {

			// Lecture d’un message WebSocket
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				log.Println("Erreur lecture:", err)
				break
			}

			// Décodage JSON du message reçu
			var incoming IncomingMsg
			if err := json.Unmarshal(msgBytes, &incoming); err != nil {
				log.Println("Erreur JSON:", err)
				continue
			}

			// Ajout des informations serveur (sécurité)
			incoming.SenderID = userID
			incoming.SenderName = name

			// Gestion selon le type de message
			switch incoming.Type {

			// Message classique
			case "message":
				// Envoi dans le channel broadcast
				// qui sera traité par HandleMessages
				broadcast <- incoming

			// Demande d’historique
			case "get_history":
				SendHistory(
					conn,
					userID,
					incoming.ReceiverID,
					incoming.Offset,
				)
			}
		}
	}
}
