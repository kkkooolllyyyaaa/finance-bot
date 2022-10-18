package messages

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type Model struct {
	sender MessageSender
}

func New(tgClient MessageSender) *Model {
	return &Model{
		sender: tgClient,
	}
}

type Message struct {
	Text   string
	UserID int64
}

func (s *Model) IncomingMessage(msg Message) error {
	payload := "Unknown command"
	if msg.Text == "/start" {
		payload = "Hello, World!"
	}
	return s.sender.SendMessage(payload, msg.UserID)
}
