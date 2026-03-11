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

// Json reçu depuis le frontend
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// représente le JSON renvoyé au client
type LoginResponse struct {
	Success bool `json:"success"`
	User    struct {
		Nickname string `json:"nickname"`
	} `json:"user"`
}
