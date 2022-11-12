package inmem_config

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/util"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/util/file"
)

type ConfigManager struct {
	data map[string]string
}

func NewConfigManager(reader file.Reader, filePath string) (*ConfigManager, error) {
	mp, err := reader.ReadToMap(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "NewConfigManager")
	}
	return &ConfigManager{
		data: mp,
	}, nil
}

func (cm *ConfigManager) currentValue(name string) string {
	return cm.data[name]
}

func (cm *ConfigManager) UpdateTimeout() int {
	defaultValue := 60
	currentValueFromConfig := cm.currentValue("timeout")

	timeout, err := util.ParseInt(currentValueFromConfig)
	if err != nil {
		return defaultValue
	}
	return timeout
}

func (cm *ConfigManager) MessagingBufferSize() int {
	defaultValue := 10
	currentValueFromConfig := cm.currentValue("messaging-buffer-size")

	messagingBufferSize, err := util.ParseInt(currentValueFromConfig)
	if err != nil {
		return defaultValue
	}
	return messagingBufferSize
}

func (cm *ConfigManager) MessagingRetryCount() int {
	defaultValue := 5
	currentValueFromConfig := cm.currentValue("messaging-retry-count")

	messagingRetryCount, err := util.ParseInt(currentValueFromConfig)
	if err != nil {
		return defaultValue
	}
	return messagingRetryCount
}

func (cm *ConfigManager) MessagingWorkersCount() int {
	defaultValue := 5
	currentValueFromConfig := cm.currentValue("messaging-workers-count")

	messagingWorkersCount, err := util.ParseInt(currentValueFromConfig)
	if err != nil {
		return defaultValue
	}
	return messagingWorkersCount
}

func (cm *ConfigManager) CurrencyRatesUpdateTick() int {
	defaultValue := 60
	currentValueFromConfig := cm.currentValue("currency-rates-update-tick")

	currencyRatesUpdateTick, err := util.ParseInt(currentValueFromConfig)
	if err != nil {
		return defaultValue
	}
	return currencyRatesUpdateTick
}
