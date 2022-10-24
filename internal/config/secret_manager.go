package config

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/reader"
)

const secretsFile = "config/secrets/secrets.yaml"

type SecretManager struct {
	data map[string]string
}

func NewSecretManager(readerModel *reader.Model) (*SecretManager, error) {
	mp, err := readerModel.ReadToMap(secretsFile)
	if err != nil {
		return nil, errors.Wrap(err, "SecretManager NewSecretManager")
	}
	return &SecretManager{
		data: mp,
	}, nil
}

func (s *SecretManager) Token() string {
	return s.data["token"]
}
