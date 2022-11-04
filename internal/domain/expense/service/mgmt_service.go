package service

import (
	"time"

	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/expense"
)

type MgmtService struct {
	repository expense.Repository
}

func NewMgmtService(repository expense.Repository) *MgmtService {
	return &MgmtService{
		repository: repository,
	}
}

func (s *MgmtService) Add(expense *entity.Expense) error {
	return s.repository.Add(expense)
}

func (s *MgmtService) GetCategoriesExpenses(userID int64, daysMinus int64) (result []*entity.CategoryExpense) {
	allByCategories := s.repository.GetAllGroupedByCategories(userID)

	now := time.Now()
	for category, expenses := range allByCategories {
		from := now.AddDate(0, 0, int(-daysMinus))
		amount := sumBetween(expenses, from, now)

		result = append(result, entity.NewCategoryExpense(amount, category))
	}

	return
}

func sumBetween(expenses []*entity.Expense, from, to time.Time) (acc float64) {
	for _, e := range expenses {
		if e.Time.Before(to) && e.Time.After(from) {
			acc += e.Amount
		}
	}
	return
}
