package entity

import "time"

type Expense struct {
	UserID   int64
	Amount   float64
	Category string
	Time     time.Time
}

type CategoryExpense struct {
	Amount   float64
	Category string
}

func NewCategoryExpense(amount float64, category string) *CategoryExpense {
	return &CategoryExpense{
		Amount:   amount,
		Category: category,
	}
}

func NewExpense(userID int64, amount float64, category string, time time.Time) *Expense {
	return &Expense{
		UserID:   userID,
		Amount:   amount,
		Category: category,
		Time:     time,
	}
}
