package service

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/common"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/command"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/util"
	"gitlab.ozon.dev/kolya_cypandin/project-base/pkg/logger"
)

type CommandService struct {
	commands map[string]command.Command
}

func NewRegistry() *CommandService {
	result := CommandService{
		commands: make(map[string]command.Command),
	}
	return &result
}

func (cs *CommandService) AddCommand(commandName string, command command.Command) {
	cs.commands[commandName] = command
}

func (cs *CommandService) Commands() map[string]command.Command {
	return cs.commands
}

var ErrCommandNotFound = errors.New("Command is not found")

func (cs *CommandService) ExecuteCommands(
	ctx context.Context,
	gotMessagesChan <-chan entity.Message,
	toSendMessages chan<- entity.Message,
) error {
	logger.Info(ctx).Msg("Start executing commands...")
	for {
		select {
		case msg := <-gotMessagesChan:
			logger.Debug(ctx).Msg("Start processing message...")
			cmd, args, err := extractArguments(msg.Text, msg.UserID)
			if err != nil {
				logger.Error(ctx).Err(err).Msg("Got error while parsing arguments:")
				writeMessage(msg, toSendMessages, errToPublicMessage(err))
				continue
			}

			logger.Debug(ctx).Msgf("Executing command=%s with args=%s", cmd, args)
			toSend, err := cs.executeCommand(cmd, args)
			if err != nil {
				logger.Warn(ctx).Err(err).Msg("Got error while executing command")
				writeMessage(msg, toSendMessages, errToPublicMessage(err))
				continue
			}

			logger.Debug(ctx).Msg("Command was successfully executed")
			writeMessage(msg, toSendMessages, toSend)
		case <-ctx.Done():
			logger.Debug(ctx).Msg("domain.command.service.ExecuteCommands context done")
			return ctx.Err()
		}
	}
}

func (cs *CommandService) executeCommand(commandName string, args []string) (result string, err error) {
	cmd, ok := cs.commands[commandName]
	if !ok {
		return result, ErrCommandNotFound
	}
	return cmd.Execute(args)
}

func writeMessage(gotMessage entity.Message, toSendMessages chan<- entity.Message, toSend string) {
	msgToSend := entity.NewMessage(toSend, gotMessage.UserID, gotMessage.MessageID)
	toSendMessages <- *msgToSend
}

func errToPublicMessage(err error) string {
	if errors.Is(err, ErrCommandNotFound) {
		return common.UnknownCommand
	}
	if errors.Is(err, common.ErrIncorrectArgument) {
		return common.CommandIncorrectFormatError
	}
	if errors.Is(err, common.ErrIncorrectArgsCount) {
		return common.CommandIncorrectArgsCountError
	}
	if errors.Is(err, util.ErrNotACommand) {
		return common.IsNotACommandError
	}
	return common.CommandExecutionError
}

func extractArguments(text string, userID int64) (cmd string, args []string, err error) {
	cmd, args, err = util.ParseCommand(text)
	if err != nil {
		return cmd, args, err
	}

	args = append([]string{
		util.FormatInt(userID),
	}, args...)

	return cmd, args, nil
}
