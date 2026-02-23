package messages

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID         int
	SenderID   int
	ReceiverID int
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Client struct {
	Conn *websocket.Conn
	Name string
	mu   sync.Mutex
}

type IncomingMsg struct {
	Type       string `json:"type"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	Offset     int    `json:"offset"`
	SenderID   int
	SenderName string
}

// IsMine ajouté pour l'expéditeur ET le destinataire
type OutgoingMsg struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	IsMine    bool   `json:"is_mine"`
}

type HistoryMsg struct {
	Sender    string `json:"sender"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	IsMine    bool   `json:"is_mine"`
}

type OnlineUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
