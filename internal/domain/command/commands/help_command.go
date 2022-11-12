package commands

import (
	"fmt"
	"strings"

	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/command"
)

type HelpCommand struct {
	commands map[string]command.Command
}

func NewHelpCommand(commands map[string]command.Command) *HelpCommand {
	return &HelpCommand{
		commands: commands,
	}
}

func (c *HelpCommand) Execute(args []string) (string, error) {
	builder := strings.Builder{}
	for k, v := range c.commands {
		msg := fmt.Sprintf("%s - %s\n", k, v.Description())
		builder.WriteString(msg)
	}
	return builder.String(), nil
}

func (c *HelpCommand) Description() string {
	return "Справочник команд"
}
