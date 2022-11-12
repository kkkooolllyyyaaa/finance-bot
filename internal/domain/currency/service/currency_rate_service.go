package service

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/client"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
	"gitlab.ozon.dev/kolya_cypandin/project-base/pkg/logger"
)

type CurrencyRatesStorage struct {
	requestService client.ApiRequestService
	currencyRates  entity.CurrencyRates
}

func New(requestService client.ApiRequestService) *CurrencyRatesStorage {
	return &CurrencyRatesStorage{
		requestService: requestService,
	}
}

func (crs *CurrencyRatesStorage) GetCurrentCurrencyRates() (currencyRates entity.CurrencyRates, err error) {
	return crs.currencyRates, nil
}

func (crs *CurrencyRatesStorage) ListenForUpdates(ctx context.Context, ticker *time.Ticker) error {
	logger.Info(ctx).Msg("Initializing currency rates...")
	err := crs.updateCurrencyRates()
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Got error while currency rates initialization")
		return err
	}

	logger.Info(ctx).Msg("Listening for currency rates updates...")
	for {
		select {
		case <-ctx.Done():
			logger.Debug(ctx).Msg("domain.currency.service.ListenForUpdates context done")
			return ctx.Err()
		case <-ticker.C:
			logger.Debug(ctx).Msg("Regular tick, will update currency rates...")
			err = crs.updateCurrencyRates()
			if err != nil {
				logger.Error(ctx).Err(err).Msg("Got error while updating currency rates, will be ignored")
			}
			logger.Debug(ctx).Msg("Currency rates successfully updated after regular tick")
		}
	}
}

func (crs *CurrencyRatesStorage) updateCurrencyRates() (err error) {
	responseBody, err := crs.requestService.Request()
	if err != nil {
		return err
	}

	rates := entity.CurrencyRates{}
	err = json.NewDecoder(responseBody).Decode(&rates)
	if err != nil {
		return err
	}

	crs.currencyRates = rates
	return nil
}
