package command

type Command interface {
	Execute([]string) (string, error)
	Description() string
}
