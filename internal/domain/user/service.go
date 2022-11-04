package user

import "gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"

type Service interface {
	CreateUserIfNotExist(userID int64) error
	SetUserCurrency(userID int64, currency entity.Currency) error
	GetUserCurrency(userID int64) (entity.Currency, error)
}
