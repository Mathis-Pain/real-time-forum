package auth

import "time"

type Session struct {
	ID        int
	UserID    int
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
	UserAgent string
	IP        string
}
