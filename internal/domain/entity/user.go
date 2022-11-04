package entity

import "time"

type User struct {
	UserID      int64
	Currency    Currency
	CreatedTime time.Time
}

func NewUser(userID int64) *User {
	return &User{
		UserID:      userID,
		Currency:    RUB,
		CreatedTime: time.Now(),
	}
}
