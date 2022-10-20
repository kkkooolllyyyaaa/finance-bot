package commands

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/common"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/expense"
	service "gitlab.ozon.dev/kolya_cypandin/project-base/internal/service/expense"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/util"
	"strings"
)

type ReportCommand struct {
	service *service.MgmtService
}

func NewReportCommand(service *service.MgmtService) *ReportCommand {
	return &ReportCommand{
		service: service,
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

	categoriesExpenses := c.service.GetCategoriesExpenses(userID, days)
	var categoriesAmount float64 = 0
	for _, ce := range categoriesExpenses {
		categoriesAmount += ce.Amount
	}

	return reportCommandMessage(categoriesExpenses, categoriesAmount), nil
}

func (c *ReportCommand) Description() string {
	builder := strings.Builder{}
	builder.WriteString("Вывести отчёт\n")
	builder.WriteString("\t\tФормат: /report [days]\n")
	builder.WriteString("\t\tПример: /report 7")
	return builder.String()
}

func reportCommandMessage(categoriesExpenses []expense.CategoryExpense, categoriesAmount float64) string {
	builder := strings.Builder{}

	builder.WriteString("Всего потрачено: ")
	categoriesAmountStr := util.FormatFloat(categoriesAmount)
	builder.WriteString(categoriesAmountStr)
	builder.WriteString("\n")

	builder.WriteString("Всего категорий: ")
	countStr := util.FormatInt(int64(len(categoriesExpenses)))
	builder.WriteString(countStr)
	builder.WriteString("\n")

	for _, ce := range categoriesExpenses {
		builder.WriteString("\t\tКатегория: ")
		builder.WriteString(ce.Category)

		builder.WriteString(" Сумма: ")
		amountStr := util.FormatFloat(ce.Amount)
		builder.WriteString(amountStr)
		builder.WriteString("\n")
	}
	return builder.String()
}

func extractReportCommandArguments(args []string) (userID, days int64, err error) {
	if len(args) != 2 {
		return userID, days, common.ErrIncorrectArgsCount
	}

	userID, err = util.ParseInt64(args[0])
	if err != nil {
		return userID, days, common.ErrIncorrectUserID
	}

	days, err = util.ParseInt64(args[1])
	ok := util.ValidateNatural(days)
	if err != nil || !ok {
		return userID, days, errIncorrectDays
	}

	return userID, days, nil
}
