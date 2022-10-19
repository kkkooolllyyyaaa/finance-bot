package commands

import (
	util2 "gitlab.ozon.dev/kolya_cypandin/project-base/internal/util"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/common"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/expense"
	service "gitlab.ozon.dev/kolya_cypandin/project-base/internal/service/expense"
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
	if len(args) != 2 {
		return result, common.ErrIncorrectArgsCount
	}

	userID, err := util2.ParseInt64(args[0])
	if err != nil {
		return result, common.ErrIncorrectUserID
	}

	days, err := util2.ParseInt64(args[1])
	if err != nil {
		return result, errIncorrectDays
	}
	ok := util2.ValidateNatural(days)
	if !ok {
		return result, errIncorrectDays
	}

	categoriesExpenses := c.service.GetCategoriesExpenses(userID, days)

	return reportCommandMessage(categoriesExpenses), nil
}

func (c *ReportCommand) Description() string {
	builder := strings.Builder{}
	builder.WriteString("Вывести отчёт\n")
	builder.WriteString("\t\tФормат: /report [days]\n")
	builder.WriteString("\t\tПример: /report 7")
	return builder.String()
}

func reportCommandMessage(categoriesExpenses []expense.CategoryExpense) string {
	builder := strings.Builder{}

	builder.WriteString("Всего категорий: ")
	countStr := util2.FormatInt(int64(len(categoriesExpenses)))
	builder.WriteString(countStr)
	builder.WriteString("\n")

	for _, ce := range categoriesExpenses {
		builder.WriteString("\t\tКатегория: ")
		builder.WriteString(ce.Category)

		builder.WriteString(" Сумма: ")
		amountStr := util2.FormatFloat(ce.Amount)
		builder.WriteString(amountStr)
		builder.WriteString("\n")
	}
	return builder.String()
}
