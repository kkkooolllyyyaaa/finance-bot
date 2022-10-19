package commands

import (
	util2 "gitlab.ozon.dev/kolya_cypandin/project-base/internal/util"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/common"
	expenseModel "gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/expense"
	service "gitlab.ozon.dev/kolya_cypandin/project-base/internal/service/expense"
)

type AddExpenseCommand struct {
	service *service.MgmtService
}

func NewAddExpenseCommand(service *service.MgmtService) *AddExpenseCommand {
	return &AddExpenseCommand{
		service: service,
	}
}

var errIncorrectAmount = errors.Wrap(
	common.ErrIncorrectArgument,
	"Incorrect amount, must be non-negative number",
)

func (c *AddExpenseCommand) Execute(args []string) (result string, err error) {
	if len(args) != 3 && len(args) != 4 {
		return result, common.ErrIncorrectArgsCount
	}

	userID, err := util2.ParseInt64(args[0])
	if err != nil {
		return result, common.ErrIncorrectUserID
	}
	amount, err := util2.ParseFloat64(args[1])
	if err != nil {
		return result, errIncorrectAmount
	}
	if amount < 0 {
		return result, errIncorrectAmount
	}
	category := args[2]
	date := time.Now()

	e := *expenseModel.NewExpense(userID, amount, category, date)
	err = c.service.Add(e)
	if err != nil {
		return result, errors.Wrap(err, "Can't add expense")
	}

	return addExpenseCommandMessage(e), nil
}

func (c *AddExpenseCommand) Description() string {
	builder := strings.Builder{}
	builder.WriteString("Добавить трату\n")
	builder.WriteString("\t\tФормат: /add [amount] [category]\n")
	builder.WriteString("\t\tПример: /add 200.0 такси")
	return builder.String()
}

func addExpenseCommandMessage(e expenseModel.Expense) string {
	builder := strings.Builder{}
	builder.WriteString("Трата успешно добавлена:\n")

	builder.WriteString("\t\tСумма: ")
	amountStr := util2.FormatFloat(e.Amount)
	builder.WriteString(amountStr)
	builder.WriteString("\n")

	builder.WriteString("\t\tКатегория: ")
	builder.WriteString(e.Category)
	builder.WriteString("\n")

	return builder.String()
}
