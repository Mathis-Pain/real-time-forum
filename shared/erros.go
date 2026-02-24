package shared

import (
	"encoding/json"
	"net/http"
)

// Fonction qui envoie message d'erreur côté client (console)
func RespondError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
