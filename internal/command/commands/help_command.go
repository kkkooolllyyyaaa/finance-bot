package commands

import (
	"strings"

	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/command"
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
		builder.WriteString(k)
		builder.WriteString(" - ")
		builder.WriteString(v.Description())
		builder.WriteString("\n")
	}
	return builder.String(), nil
}

func (c *HelpCommand) Description() string {
	return "Справочник команд"
}
