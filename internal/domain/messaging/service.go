package messaging

import (
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
)

type Messaging interface {
	SendMessages(<-chan entity.Message) error
	ListenForMessages() (<-chan entity.Message, error)
}
