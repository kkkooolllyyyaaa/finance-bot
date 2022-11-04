package user

import "gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"

type Repository interface {
	Create(user *entity.User) error
	Get(userID int64) (*entity.User, error)
	UpdateCurrency(userID int64, currency entity.Currency) error
}
