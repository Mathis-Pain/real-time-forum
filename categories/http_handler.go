package categories

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func CategoriesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cats, err := GetAllCategories(db)
		if err != nil {
			http.Error(w, "Erreur récupération catégories", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cats)
	}
}
