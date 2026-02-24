package posts

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/shared"
	"strconv"
)

func PostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newPost Post
		var response PostResponse
		switch r.Method {
		case "POST":
			if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				shared.RespondError(w, http.StatusBadRequest, err)
				return
			}
			// Récupérer l'auteur grâce au cookie

			authorId, err := shared.GetUserIdByCookie(r, db)
			if err != nil {
				log.Println("<NewPostHandler> Error cannot get userId: ", err)
				shared.RespondError(w, http.StatusInternalServerError, err)
				return
			}

			newPost.AuthorID = authorId

			// -- gestion d'erreurs --
			// ajouter verif post existant

			if err := IsValidFormat(newPost); err != nil {
				log.Println("<NewPostHandler> Error invalid post: ", err)
				shared.RespondError(w, http.StatusBadRequest, err)
				return
			}

			// -- sauvegarde db --

			if err := SavePost(db, &newPost); err != nil {
				log.Println("<NewPostHandler> Error cannot save post: ", err)
				shared.RespondError(w, http.StatusInternalServerError, err)
				return
			}

			response = PostResponse{
				ResponseNotif: "Post created",
				Post:          newPost,
			}

			log.Println("New post saved into db")
		case "GET":
			referer := r.Header.Get("Referer")
			if referer == "" {
				referer = "/"
			}

			//reception id post
			idStr := r.URL.Query().Get("id")

			postId, err := strconv.Atoi(idStr)
			if err != nil {
				log.Printf("Erreur conversion ou id invalide: %v, retour à la page précédente", err)
				shared.RespondError(w, http.StatusInternalServerError, err)
				http.Redirect(w, r, referer, http.StatusSeeOther)
				return
			}

			if postId == 0 {
				allPosts, err := GetAllPosts(db)
				if err != nil {
					log.Printf("Erreur recuperation allPosts: %v", err)
					shared.RespondError(w, http.StatusInternalServerError, err)
				}
				response = PostResponse{
					ResponseNotif: "Sendind all posts",
					AllPosts:      allPosts,
				}
			} else {
				// récuperer toutes les données du post

				currentPost, currentCmts, currentAuthor, err := GetPostData(db, postId)
				if err != nil {
					log.Printf("Erreur récupération des données du post: %v", err)
					shared.RespondError(w, http.StatusInternalServerError, err)
					return
				}

				//renvoyer les données à afficher
				response = PostResponse{
					ResponseNotif: "Sending post",
					Post:          currentPost,
					Comments:      currentCmts,
					Author:        currentAuthor,
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

}
