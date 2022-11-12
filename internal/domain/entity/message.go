package entity

type Message struct {
	Text      string
	UserID    int64
	MessageID int
}

func NewMessage(text string, userID int64, messageID int) *Message {
	return &Message{
		Text:      text,
		UserID:    userID,
		MessageID: messageID,
	}
}
