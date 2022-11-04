package repository

import (
	"fmt"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/common"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
)

type InMemRepository struct {
	users map[int64]*entity.User
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		users: make(map[int64]*entity.User, 0),
	}
}

func (r *InMemRepository) Create(user *entity.User) error {
	_, ok := r.users[user.UserID]
	if !ok {
		r.users[user.UserID] = user
		return nil
	}
	msg := fmt.Sprintf("User with id=%d already exists", user.UserID)
	return errors.Wrap(common.EntityAlreadyExists, msg)
}

func (r *InMemRepository) Get(userID int64) (*entity.User, error) {
	existing, ok := r.users[userID]
	if !ok {
		msg := fmt.Sprintf("User with id=%d isn't exist", userID)
		return nil, errors.Wrap(common.EntityNotFound, msg)
	}
	return existing, nil
}

func (r *InMemRepository) UpdateCurrency(userID int64, currency entity.Currency) error {
	existing, err := r.Get(userID)
	if err != nil {
		return err
	}
	existing.Currency = currency
	return nil
}
