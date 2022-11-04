package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/common"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/currency"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/expense"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/user"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/util"
)

type AddExpenseCommand struct {
	expenseService  expense.Service
	userService     user.Service
	currencyService currency.Service
}

func NewAddExpenseCommand(
	service expense.Service,
	userService user.Service,
	currencyService currency.Service,
) *AddExpenseCommand {
	return &AddExpenseCommand{
		expenseService:  service,
		userService:     userService,
		currencyService: currencyService,
	}
}

var errIncorrectAmount = errors.Wrap(
	common.ErrIncorrectArgument,
	"Incorrect amount, must be non-negative number",
)

func (c *AddExpenseCommand) Execute(args []string) (result string, err error) {
	userID, amount, category, err := extractAddCommandArguments(args)
	if err != nil {
		return result, errors.Wrap(err, "Error while extracting arguments")
	}

	userCurrency, currencySign, err := c.getCurrencyAndSign(userID)
	if err != nil {
		return result, err
	}

	currencyRates, err := c.currencyService.GetCurrentCurrencyRates()
	if err != nil {
		return result, err
	}

	amountRub, err := currencyRates.ConvertFromTo(amount, userCurrency, entity.RUB)
	if err != nil {
		return result, err
	}
	expenseToAdd := entity.NewExpense(userID, amountRub, category, time.Now())
	err = c.expenseService.Add(expenseToAdd)
	if err != nil {
		return result, err
	}

	return addExpenseCommandMessage(amount, expenseToAdd.Category, currencySign), nil
}

func (c *AddExpenseCommand) Description() string {
	builder := strings.Builder{}
	builder.WriteString("Добавить трату\n")
	builder.WriteString("\t\tФормат: /add [amount] [category]\n")
	builder.WriteString("\t\tПример: /add 200.0 такси")
	return builder.String()
}

func addExpenseCommandMessage(amount float64, category, currencySign string) string {
	builder := strings.Builder{}
	builder.WriteString("Трата успешно добавлена:\n")

	msg := fmt.Sprintf("\t\tСумма: %.2f %s\n", amount, currencySign)
	builder.WriteString(msg)

	msg = fmt.Sprintf("\t\tКатегория: %s\n", category)
	builder.WriteString(msg)

	return builder.String()
}

func (c *AddExpenseCommand) getCurrencyAndSign(userID int64) (currency entity.Currency, sign string, err error) {
	userCurrency, err := c.userService.GetUserCurrency(userID)
	if err != nil {
		return currency, sign, errors.Wrap(err, "Unable to get user currency")
	}
	currencySign, ok := entity.GetCurrencySign(userCurrency)
	if !ok {
		msg := fmt.Sprintf("Sign for currency=%s isn't exist", userCurrency)
		return currency, sign, errors.Wrap(common.EntityNotFound, msg)
	}

	return userCurrency, currencySign, nil
}

func extractAddCommandArguments(args []string) (userID int64, amount float64, category string, err error) {
	if len(args) != 3 && len(args) != 4 {
		return userID, amount, category, common.ErrIncorrectArgsCount
	}
	userID, err = util.ParseInt64(args[0])
	if err != nil {
		return userID, amount, category, common.ErrIncorrectUserID
	}
	amount, err = util.ParseFloat64(args[1])
	if err != nil {
		return userID, amount, category, errIncorrectAmount
	}
	if amount < 0 {
		return userID, amount, category, errIncorrectAmount
	}
	return userID, amount, args[2], err
}
