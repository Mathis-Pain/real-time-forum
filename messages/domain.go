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

// ajout pour typing
type DirectMessage struct {
	To   string
	Data []byte
}
type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Send       chan DirectMessage // ajout pour typing
}

var HubInstance = Hub{
	Clients:    make(map[string]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Send:       make(chan DirectMessage, 256), // ajout pour typing
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
