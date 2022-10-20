package commands

import (
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/common"
	"strings"
)

type StartCommand struct{}

func NewStartCommand() *StartCommand {
	return &StartCommand{}
}

func (c *StartCommand) Execute(args []string) (result string, err error) {
	return common.StartMessage, nil
}

func (c *StartCommand) Description() string {
	builder := strings.Builder{}
	builder.WriteString("Начало работы бота\n")
	return builder.String()
}
