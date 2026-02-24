package comments

import "time"

type Comment struct {
	ID        int
	PostID    int    `json:"postid"`
	AuthorID  int    `json:"authorid"`
	Content   string `json:"content"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentResponse struct {
	ID       int
	PostID   int    `json:"postid"`
	AuthorID int    `json:"authorid"`
	UserName string `json:"authorname"`
	Content  string `json:"content"`
}
