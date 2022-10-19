package command

import (
	"github.com/pkg/errors"
)

type Registry struct {
	commands map[string]Command
}

func NewRegistry() *Registry {
	result := Registry{
		commands: make(map[string]Command),
	}
	return &result
}

var ErrCommandNotFound = errors.New("Command is not found")

func (r *Registry) Execute(commandName string, args []string) (result string, err error) {
	cmd, ok := r.commands[commandName]
	if !ok {
		return result, ErrCommandNotFound
	}
	return cmd.Execute(args)
}

func (r *Registry) AddCommand(commandName string, command Command) {
	r.commands[commandName] = command
}

func (r *Registry) Commands() map[string]Command {
	return r.commands
}
