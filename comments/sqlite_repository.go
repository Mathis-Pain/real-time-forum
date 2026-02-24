package comments

import (
	"database/sql"
	"fmt"
)

//IsValidMeta : exist post/author

func SaveComment(db *sql.DB, comment *Comment) error {
	query := `INSERT INTO comments (postid, authorid, content) VALUES (?, ?, ?)`
	_, err := db.Exec(query, comment.PostID, comment.AuthorID, comment.Content)
	return err
}

// --- Validité metadata ---

func ValidateCommentReferences(db *sql.DB, comment *Comment) error {
	var postID int
	err := db.QueryRow(
		`SELECT id FROM post WHERE id = ?`,
		comment.PostID,
	).Scan(&postID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("post %d does not exist", comment.PostID)
		}
		return err
	}

	// Vérifier l’auteur
	var authorID int
	err = db.QueryRow(
		`SELECT id FROM users WHERE id = ?`,
		comment.AuthorID,
	).Scan(&authorID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("author %d does not exist", comment.AuthorID)
		}
		return err
	}

	return nil
}
