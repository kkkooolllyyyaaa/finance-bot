package commands

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/common"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/currency"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/expense"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/user"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/util"
)

type ReportCommand struct {
	expenseService  expense.Service
	userService     user.Service
	currencyService currency.Service
}

func NewReportCommand(
	expenseService expense.Service,
	userService user.Service,
	currencyService currency.Service,
) *ReportCommand {
	return &ReportCommand{
		expenseService:  expenseService,
		userService:     userService,
		currencyService: currencyService,
	}
}

var errIncorrectDays = errors.Wrap(
	common.ErrIncorrectArgument,
	"Incorrect days count, should be natural number",
)

func (c *ReportCommand) Execute(args []string) (result string, err error) {
	userID, days, err := extractReportCommandArguments(args)
	if err != nil {
		return result, errors.Wrap(err, "Error while extracting arguments")
	}

	categoriesExpenses := c.expenseService.GetCategoriesExpenses(userID, days)
	var categoriesAmount float64 = 0
	for _, ce := range categoriesExpenses {
		categoriesAmount += ce.Amount
	}

	userCurrency, currencySign, err := c.getCurrencyAndSign(userID)
	if err != nil {
		return result, err
	}

	currencyRates, err := c.currencyService.GetCurrentCurrencyRates()
	if err != nil {
		return result, err
	}

	return c.enrichAndComposeMessage(
		categoriesExpenses,
		categoriesAmount,
		userCurrency,
		currencySign,
		currencyRates,
		days,
	)
}

func (c *ReportCommand) Description() string {
	builder := strings.Builder{}
	builder.WriteString("Вывести отчёт\n")
	builder.WriteString("\t\tФормат: /report [days]\n")
	builder.WriteString("\t\tПример: /report 7")
	return builder.String()
}

func (c *ReportCommand) enrichAndComposeMessage(
	categoriesExpenses []*entity.CategoryExpense,
	categoriesAmount float64,
	userCurrency entity.Currency,
	currencySign string,
	currencyRates entity.CurrencyRates,
	days int64,
) (result string, err error) {
	categoriesAmount, err = currencyRates.ConvertFromTo(categoriesAmount, entity.RUB, userCurrency)
	if err != nil {
		return result, err
	}

	builder := strings.Builder{}
	msg := fmt.Sprintf("Всего потрачено за %d дней: %.2f %s\n", days, categoriesAmount, currencySign)
	builder.WriteString(msg)

	msg = fmt.Sprintf("Всего категорий: %d\n", len(categoriesExpenses))
	builder.WriteString(msg)

	for _, ce := range categoriesExpenses {
		msg = fmt.Sprintf("\t\tКатегория: %s ", ce.Category)
		builder.WriteString(msg)

		userCurrencyAmount, err := currencyRates.ConvertFromTo(ce.Amount, entity.RUB, userCurrency)
		if err != nil {
			return result, err
		}

		msg = fmt.Sprintf("Сумма: %.2f %s\n", userCurrencyAmount, currencySign)
		builder.WriteString(msg)
	}
	return builder.String(), nil
}

func (c *ReportCommand) getCurrencyAndSign(userID int64) (currency entity.Currency, sign string, err error) {
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

func extractReportCommandArguments(args []string) (userID, days int64, err error) {
	if len(args) != 1 && len(args) != 2 {
		return userID, days, common.ErrIncorrectArgsCount
	}

	userID, err = util.ParseInt64(args[0])
	if err != nil {
		return userID, days, common.ErrIncorrectUserID
	}

	if len(args) == 1 {
		return userID, 7, nil
	}

	days, err = util.ParseInt64(args[1])
	ok := util.ValidateNatural(days)
	if err != nil || !ok {
		return userID, days, errIncorrectDays
	}

	return userID, days, nil
}
