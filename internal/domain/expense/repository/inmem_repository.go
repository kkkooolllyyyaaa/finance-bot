package repository

import (
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
)

type InMemRepository struct {
	expenses map[int64][]*entity.Expense
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		expenses: make(map[int64][]*entity.Expense, 0),
	}
}

func (r *InMemRepository) Add(expense *entity.Expense) error {
	r.expenses[expense.UserID] = append(r.expenses[expense.UserID], expense)
	return nil
}

func (r *InMemRepository) GetAllOfUser(userID int64) []*entity.Expense {
	return r.expenses[userID]
}

func (r *InMemRepository) GetAllByCategoryOfUser(userID int64, category string) (acc []*entity.Expense) {
	all := r.GetAllOfUser(userID)
	for _, e := range all {
		if e.Category == category {
			acc = append(acc, e)
		}
	}
	return
}

func (r *InMemRepository) GetAllGroupedByCategories(userID int64) map[string][]*entity.Expense {
	categoryToExpenses := make(map[string][]*entity.Expense)

	all := r.GetAllOfUser(userID)
	for _, e := range all {
		categoryToExpenses[e.Category] = append(categoryToExpenses[e.Category], e)
	}

	return categoryToExpenses
}
