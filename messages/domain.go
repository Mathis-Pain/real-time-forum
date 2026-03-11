package messages

import (
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
type lastUserMessage struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	CreatedAt time.Time `json:"created_at"`
}

type HistoryMessage struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Client struct {
	UserName string
	Conn     *websocket.Conn
	Send     chan []byte
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
}

var HubInstance = Hub{
	Clients:    make(map[string]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

type IncomingMessage struct {
	Type    string `json:"type"`    // "private_message"
	To      string `json:"to"`      // receiver username
	Content string `json:"content"` // message text
}

type OutgoingMessage struct {
	Type      string    `json:"type"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
