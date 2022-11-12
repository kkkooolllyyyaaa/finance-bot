package expense

import (
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
)

type Service interface {
	GetCategoriesExpenses(userID int64, daysMinus int64) (result []*entity.CategoryExpense)
	Add(expense *entity.Expense) error
}
