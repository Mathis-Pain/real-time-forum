package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"real-time-forum/auth"
	"real-time-forum/comments"
	"real-time-forum/messages"
	"real-time-forum/posts"
	"real-time-forum/users"
)

// Commande de lancement depuis la racine : go run cmd/main.go

func Server(port string, db *sql.DB) {
	mux := http.NewServeMux() // Création d'un serveur mux vide.

	// routes publiques
	mux.HandleFunc("/auth/login", auth.LoginHandler(db))
	mux.HandleFunc("/logout", auth.LogoutHandler(db))

	// routes protégées
	mux.Handle("/auth/me", auth.AuthMiddleware(db)(http.HandlerFunc(auth.CurrentUserHandler(db))))
	mux.Handle("/post", auth.AuthMiddleware(db)(http.HandlerFunc(posts.PostHandler(db))))
	mux.Handle("/comment", auth.AuthMiddleware(db)(http.HandlerFunc(comments.NewCommentHandler(db))))
	mux.HandleFunc("/ws", messages.WsHandler) // Pour le websocket.
	// mux.Handle("/ws", auth.AuthMiddleware(db)(http.HandlerFunc(messages.WsHandler)))

	userRepo := &users.Repository{DB: db}
	userHandler := &users.Handler{Repo: userRepo}
	mux.Handle("/online-users", auth.AuthMiddleware(db)(http.HandlerFunc(userHandler.OnlineUsersHandler)))
	mux.Handle("/messages/history", auth.AuthMiddleware(db)(http.HandlerFunc(messages.GetHistoryHandler)))

	// Quand l'utilisateur arrive, affiche mainPage.
	mux.HandleFunc("/", users.MainPage)
	// servir les fichiers static
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", fs))
	// Lancement serveur Go
	fmt.Printf("Serveur démarré sur le port %s \n", port)
	fmt.Printf("Page d'accès : http://localhost:8080/ \n")
	err := http.ListenAndServe(port, mux) // A laisser à la fin, élément bloquant la lecture des instructions suivantes.
	if err != nil {
		fmt.Printf("Error lauching servor : %v \n", err)
	}
}
