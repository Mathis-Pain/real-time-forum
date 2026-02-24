package posts

import (
	"database/sql"
	"fmt"
	"real-time-forum/comments"
	"real-time-forum/users"
)

// --- Fonctions Sauvegarde de Data

func SavePost(db *sql.DB, post *Post) error {
	query := `INSERT INTO post (title, content, authorid) VALUES (?, ?, ?)`
	_, err := db.Exec(query, post.Title, post.Content, post.AuthorID)
	SavePostCategories(db, post)
	return err
}

func SavePostCategories(db *sql.DB, post *Post) error {
	postWithId, err := GetPostID(db, post)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		`INSERT INTO post_categories (postID, categoryID) VALUES(?, ?)`,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, catId := range post.CategoryIDs {
		_, err := stmt.Exec(postWithId.ID, catId)
		fmt.Println(postWithId.ID, catId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// --- Fonctions Récupération de Data ---

func GetPostID(db *sql.DB, post *Post) (*Post, error) {
	query := `SELECT id FROM post WHERE title = ? AND content = ? AND authorid = ?`
	row := db.QueryRow(query, post.Title, post.Content, post.AuthorID)
	err := row.Scan(&post.ID)
	if err != nil {
		return &Post{}, err
	}
	return post, nil
}

func GetAuthorInfo(db *sql.DB, id int) (users.User, error) {
	var author users.User
	sqlQuery := `SELECT id, username FROM users WHERE id = ?`
	row := db.QueryRow(sqlQuery, id)
	if err := row.Scan(&author.ID, &author.UserName); err != nil {
		return users.User{}, err
	}

	return author, nil
}

func GetPostData(db *sql.DB, id int) (Post, []comments.CommentResponse, users.User, error) {
	var post Post
	queryPost := `SELECT id, title, content, authorid FROM post WHERE id = ?`
	row := db.QueryRow(queryPost, id)
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID)
	if err != nil {
		return Post{}, []comments.CommentResponse{}, users.User{}, err
	}

	coms := GetCommentsByPostId(db, id)

	postAuthor, err := GetAuthorInfo(db, post.AuthorID)
	if err != nil {
		return Post{}, []comments.CommentResponse{}, users.User{}, err
	}

	return post, coms, postAuthor, nil
}

func GetCommentsByPostId(db *sql.DB, postId int) []comments.CommentResponse {
	var coms []comments.CommentResponse
	var comment comments.CommentResponse
	sqlQuery := `SELECT id, authorid, content FROM comments WHERE postid = ?`
	rows, err := db.Query(sqlQuery, postId)
	if err != nil {
		return []comments.CommentResponse{}
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&comment.ID, &comment.AuthorID, &comment.Content); err != nil {
			if err == sql.ErrNoRows {
				return []comments.CommentResponse{}
			}
			return []comments.CommentResponse{}
		}

		commentAuthor, err := GetAuthorInfo(db, comment.AuthorID)
		if err != nil {
			return []comments.CommentResponse{}
		}

		comment.UserName = commentAuthor.UserName

		coms = append(coms, comment)
	}

	return coms
}

func GetAllPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT id, title, content, authorid FROM post")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.AuthorID); err != nil {
			return nil, err
		}

		// Charger les catégories du post
		catRows, err := db.Query("SELECT categoryid FROM post_categories WHERE postid = ?", p.ID)
		if err != nil {
			return nil, err
		}

		for catRows.Next() {
			var cid int
			catRows.Scan(&cid)
			p.CategoryIDs = append(p.CategoryIDs, cid)
		}
		catRows.Close()

		posts = append(posts, p)
	}

	return posts, nil
}
