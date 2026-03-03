package messages

import "time"

type Message struct {
	ID         int
	SenderID   int
	ReceiverID int
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
