package config

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/reader"
)

const configFile = "config/config.yaml"

type ConfigManager struct {
	data map[string]string
}

func NewConfigManager(readerModel *reader.Model) (*ConfigManager, error) {
	mp, err := readerModel.ReadToMap(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "ConfigManager NewConfigManager")
	}
	return &ConfigManager{
		data: mp,
	}, nil
}

func (cm *ConfigManager) UpdateTimeout() string {
	return cm.data["timeout"]
}
