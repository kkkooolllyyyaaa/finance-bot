package expense

import (
	"time"

	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/expense"
	expenseRepo "gitlab.ozon.dev/kolya_cypandin/project-base/internal/repository/expense"
)

type MgmtService struct {
	repository expenseRepo.Repository
}

func NewMgmtService(repository expenseRepo.Repository) *MgmtService {
	return &MgmtService{
		repository: repository,
	}
}

func (s *MgmtService) Add(expense *expense.Expense) error {
	return s.repository.Add(expense)
}

func (s *MgmtService) GetAll(userID int64) []*expense.Expense {
	return s.repository.GetAllOfUser(userID)
}

func (s *MgmtService) GetCategoryExpense(userID int64, category string) *expense.CategoryExpense {
	allOfCategory := s.repository.GetAllByCategoryOfUser(userID, category)

	var amount float64 = 0
	for _, e := range allOfCategory {
		amount += e.Amount
	}
	return expense.NewCategoryExpense(amount, category)
}

func (s *MgmtService) GetCategoriesExpenses(userID int64, daysMinus int64) (result []*expense.CategoryExpense) {
	allByCategories := s.repository.GetAllGroupedByCategories(userID)

	now := time.Now()
	for category, expenses := range allByCategories {
		from := now.AddDate(0, 0, int(-daysMinus))
		amount := sumBetween(expenses, from, now)

		result = append(result, expense.NewCategoryExpense(amount, category))
	}

	return
}

func sumBetween(expenses []*expense.Expense, from, to time.Time) (acc float64) {
	for _, e := range expenses {
		if e.Date.Before(to) && e.Date.After(from) {
			acc += e.Amount
		}
	}
	return
}
