package posts

import (
	"real-time-forum/comments"
	"real-time-forum/users"
	"time"
)

type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	AuthorID    int    `json:"authorid"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CategoryIDs []int `json:"category_ids"`
}

type PostResponse struct {
	ResponseNotif string                     `json:"notif"`
	AllPosts      []Post                     `json:"allposts"`
	Post          Post                       `json:"post"`
	Comments      []comments.CommentResponse `json:"comments"`
	Author        users.User                 `json:"author"`
}
