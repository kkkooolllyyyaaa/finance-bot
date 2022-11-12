package commands

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/common"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/user"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/util"
)

type SetCurrencyCommand struct {
	userService user.Service
}

func NewSetCurrencyCommand(userService user.Service) *SetCurrencyCommand {
	return &SetCurrencyCommand{
		userService: userService,
	}
}

func (c *SetCurrencyCommand) Execute(args []string) (result string, err error) {
	userID, currency, err := extractSetCurrencyCommand(args)
	if err != nil {
		return result, errors.Wrap(err, "Error while extracting arguments")
	}

	err = c.userService.SetUserCurrency(userID, currency)
	if err != nil {
		return result, err
	}

	return setCurrencyCommandMessage(currency), nil
}

func (c *SetCurrencyCommand) Description() string {
	builder := strings.Builder{}
	builder.WriteString("Выбрать отображаемую валюту\n")
	builder.WriteString("\t\tФормат: /currency [currency]\n")
	builder.WriteString("\t\t\t\tВозможные значения currency [ RUB, USD, EUR, CNY ]")
	builder.WriteString("\t\tПример: /currency USD")
	return builder.String()
}

func setCurrencyCommandMessage(currency entity.Currency) string {
	builder := strings.Builder{}
	msg := fmt.Sprintf("Валюта %s успешно выбрана!", currency)
	builder.WriteString(msg)
	return builder.String()
}

func extractSetCurrencyCommand(args []string) (userID int64, currency entity.Currency, err error) {
	if len(args) != 2 {
		return userID, currency, common.ErrIncorrectArgsCount
	}

	userID, err = util.ParseInt64(args[0])
	if err != nil {
		return userID, currency, common.ErrIncorrectUserID
	}

	currency, ok := entity.GetCurrency(args[1])
	if !ok {
		return userID, currency, errors.New("Incorrect currency")
	}

	return userID, currency, err
}
