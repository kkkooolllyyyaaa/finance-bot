package command

import (
	"context"

	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
)

type Service interface {
	ExecuteCommands(
		ctx context.Context,
		gotMessagesChan <-chan entity.Message,
		toSendMessages chan<- entity.Message,
	) error

	AddCommand(commandName string, command Command)

	Commands() map[string]Command
}
