package expense

import (
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
)

type Repository interface {
	Add(*entity.Expense) error
	GetAllOfUser(int64) []*entity.Expense
	GetAllByCategoryOfUser(int64, string) []*entity.Expense
	GetAllGroupedByCategories(int64) map[string][]*entity.Expense
}
