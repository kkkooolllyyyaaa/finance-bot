package expense

import (
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/expense"
)

type InMemRepository struct {
	expenses map[int64][]expense.Expense
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		expenses: make(map[int64][]expense.Expense, 0),
	}
}

func (r *InMemRepository) Add(expense *expense.Expense) error {
	r.expenses[expense.UserID] = append(r.expenses[expense.UserID], *expense)
	return nil
}

func (r *InMemRepository) GetAllOfUser(userID int64) []expense.Expense {
	return r.expenses[userID]
}

func (r *InMemRepository) GetAllByCategoryOfUser(userID int64, category string) (acc []expense.Expense) {
	all := r.GetAllOfUser(userID)
	for _, e := range all {
		if e.Category == category {
			acc = append(acc, e)
		}
	}
	return
}

func (r *InMemRepository) GetAllGroupedByCategories(userID int64) map[string][]expense.Expense {
	categoryToExpenses := make(map[string][]expense.Expense)

	all := r.GetAllOfUser(userID)
	for _, e := range all {
		categoryToExpenses[e.Category] = append(categoryToExpenses[e.Category], e)
	}

	return categoryToExpenses
}
