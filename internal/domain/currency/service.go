package currency

import (
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
)

type Service interface {
	GetCurrentCurrencyRates() (entity.CurrencyRates, error)
}
