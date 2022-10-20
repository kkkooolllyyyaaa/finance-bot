package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/common"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/config"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/command"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/messages"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/util"
	"log"
)

type Tg struct {
	tg            *tgbotapi.BotAPI
	cmdRegistry   *command.Registry
	configManager *config.ConfigManager
}

func New(token string, registry *command.Registry, manager *config.ConfigManager) (*Tg, error) {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "NewBotApi")
	}

	return &Tg{
		tg:            client,
		cmdRegistry:   registry,
		configManager: manager,
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
	u.Timeout, _ = util.ParseInt(c.configManager.UpdateTimeout())
	updates := c.tg.GetUpdatesChan(u)

	log.Println("Listening for messages...")
	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text
		userID := update.Message.From.ID

		log.Println("Trying to execute got command...")
		toSend := c.tryToExecute(text, userID)

		log.Println("Sending message...")
		if err := msgModel.Send(
			messages.Message{
				Text:   toSend,
				UserID: userID,
			},
		); err != nil {
			log.Println("Error sending message:", err)
			continue
		}
	}
}

func (c *Tg) tryToExecute(text string, userID int64) string {
	log.Printf("[%d] %s", userID, text)
	cmd, args, err := util.ParseCommand(text)
	if err != nil {
		log.Println("error parsing command", err)
		return common.IsNotACommandError
	}

	args = append([]string{
		util.FormatInt(userID),
	}, args...)

	log.Printf("Executing command=%s with args=%s", cmd, args)
	result, err := c.cmdRegistry.Execute(cmd, args)
	if err != nil {
		log.Printf("error executing command cmd=%s, args=%s, err=%s", cmd, args, err)
		return errToPublicMessage(err)
	}

	return result
}

func errToPublicMessage(err error) string {
	if errors.Is(err, command.ErrCommandNotFound) {
		return common.UnknownCommand
	}
	if errors.Is(err, common.ErrIncorrectArgument) {
		return common.CommandIncorrectFormatError
	}
	if errors.Is(err, common.ErrIncorrectArgsCount) {
		return common.CommandIncorrectArgsCountError
	}
	return common.CommandExecutionError
}
