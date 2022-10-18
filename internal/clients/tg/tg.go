package tg

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/messages"
)

type Tg struct {
	tg *tgbotapi.BotAPI
}

func New(token string) (*Tg, error) {
	client, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, errors.Wrap(err, "NewBotApi")
	}

	return &Tg{
		tg: client,
	}, nil
}
func (c *Tg) SendMessage(text string, userID int64) error {
	message := tgbotapi.NewMessage(userID, text)

	_, err := c.tg.Send(message)

	if err != nil {
		return errors.Wrap(err, "tg.Send")
	}
	return nil
}

func (c *Tg) ListenUpdates(msgModel *messages.Model) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.tg.GetUpdatesChan(u)

	log.Println("listening for messages")

	for update := range updates {
		if update.Message == nil {
			return
		}

		text := update.Message.Text
		log.Printf("[%s] %s", update.Message.From.UserName, text)

		err := msgModel.IncomingMessage(messages.Message{
			Text:   text,
			UserID: update.Message.From.ID,
		})

		if err != nil {
			log.Println("error processing message:", err)
		}
	}
}
