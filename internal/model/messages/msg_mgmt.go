package messages

type MessageSender interface {
	SendMessage(text string, userID int64) error
}

type Model struct {
	sender MessageSender
}

func New(sender MessageSender) *Model {
	return &Model{
		sender: sender,
	}
}

type Message struct {
	Text   string
	UserID int64
}

func (m *Model) Send(msg *Message) error {
	return m.sender.SendMessage(msg.Text, msg.UserID)
}
