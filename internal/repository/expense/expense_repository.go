package expense

import (
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/expense"
)

type Repository interface {
	Add(*expense.Expense) error
	GetAllOfUser(int64) []expense.Expense
	GetAllByCategoryOfUser(int64, string) []expense.Expense
	GetAllGroupedByCategories(int64) map[string][]expense.Expense
}
