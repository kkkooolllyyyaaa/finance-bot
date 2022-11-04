package service

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/common"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/user"
)

type MgmtService struct {
	repository user.Repository
}

func NewMgmtService(repo user.Repository) *MgmtService {
	return &MgmtService{repository: repo}
}

func (s *MgmtService) CreateUserIfNotExist(userID int64) error {
	newUser := entity.NewUser(userID)
	err := s.repository.Create(newUser)

	if errors.Is(err, common.EntityAlreadyExists) {
		return nil
	} else {
		return err
	}
}

func (s *MgmtService) SetUserCurrency(userID int64, currency entity.Currency) error {
	return s.repository.UpdateCurrency(userID, currency)
}

func (s *MgmtService) GetUserCurrency(userID int64) (currency entity.Currency, err error) {
	existing, err := s.repository.Get(userID)
	if err != nil {
		return currency, err
	}

	return existing.Currency, nil
}
