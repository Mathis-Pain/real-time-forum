package comments

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/shared"
)

func NewCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//récup data

		var newComment Comment

		if r.Method == "POST" {
			if err := json.NewDecoder(r.Body).Decode(&newComment); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		//authorID

		authorId, err := shared.GetUserIdByCookie(r, db)
		if err != nil {
			log.Println("<NewCommentHandler> Error cannot get cookie: ", err)
			shared.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		newComment.AuthorID = authorId

		//check data
		if err := IsValidFormat(newComment); err != nil {
			log.Println("<NewCommentHandler> Error invalid comment: ", err)
			shared.RespondError(w, http.StatusBadRequest, err)
			return
		}

		if err := ValidateCommentReferences(db, &newComment); err != nil {
			log.Println("<NewCommentHandler> Error invalid comment references: ", err)
			shared.RespondError(w, http.StatusBadRequest, err)
			return
		}

		//save data
		if err := SaveComment(db, &newComment); err != nil {
			log.Println("<NewCommentHandler> Error cannot save comment: ", err)
			shared.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		// --- notification comment ajouté ---
		log.Println("New comment saved into db")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "comment created",
		})
	}
}
