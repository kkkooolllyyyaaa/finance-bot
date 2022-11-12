package service

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/config"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/user"
	"gitlab.ozon.dev/kolya_cypandin/project-base/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type TgMessaging struct {
	tgApi         *tgbotapi.BotAPI
	userService   user.Service
	configManager config.ConfigManager
	retryCount    int
}

func NewTgMessaging(
	token string,
	userService user.Service,
	configManager config.ConfigManager,
) (*TgMessaging, error) {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "NewBotApi")
	}

	return &TgMessaging{
		tgApi:         client,
		userService:   userService,
		configManager: configManager,
		retryCount:    configManager.MessagingRetryCount(),
	}, nil
}

func (tg *TgMessaging) SendMessages(ctx context.Context, erg *errgroup.Group, messagesChannel <-chan entity.Message) error {
	logger.Info(ctx).Msg("Start sending messages...")
	for {
		select {
		case message := <-messagesChannel:
			logger.Debug(ctx).Msgf("Sending message to user=%d", message.UserID)

			toSend := tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:           message.UserID,
					ReplyToMessageID: message.MessageID,
				},
				Text:                  message.Text,
				DisableWebPagePreview: false,
			}
			erg.Go(func() error {
				return tg.retryableSend(toSend, 0)
			})
		case <-ctx.Done():
			logger.Debug(ctx).Msg("domain.messaging.service.SendMessages context done")
			return ctx.Err()
		}
	}
}

func (tg *TgMessaging) ListenForMessages(ctx context.Context, messageChan chan<- entity.Message) error {
	u := tgbotapi.NewUpdate(0)
	updateTimeout := tg.configManager.UpdateTimeout()
	u.Timeout = updateTimeout
	updates := tg.tgApi.GetUpdatesChan(u)

	logger.Info(ctx).Msg("Listening for messages...")
	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			_ = tg.userService.CreateUserIfNotExist(update.Message.From.ID)
			message := entity.NewMessage(update.Message.Text, update.Message.From.ID, update.Message.MessageID)
			logger.Debug(ctx).Msgf("[%d]: %s", message.UserID, message.Text)
			messageChan <- *message
		case <-ctx.Done():
			logger.Debug(ctx).Msg("domain.messaging.service.ListenForMessages context done")
			return ctx.Err()
		}
	}
}

func (tg *TgMessaging) retryableSend(toSend tgbotapi.Chattable, cnt int) error {
	_, err := tg.tgApi.Send(toSend)
	if err != nil && cnt < tg.retryCount {
		return tg.retryableSend(toSend, cnt+1)
	} else if err != nil {
		return err
	}
	return nil
}
