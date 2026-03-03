package users

import (
	"encoding/json"
	"net/http"
)

// Exporté : Handler
type Handler struct {
	Repo *Repository
}

// Exporté : OnlineUsersHandler
func (h *Handler) OnlineUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.Repo.GetOnlineUsers()
	if err != nil {
		http.Error(w, "Erreur serveur", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
