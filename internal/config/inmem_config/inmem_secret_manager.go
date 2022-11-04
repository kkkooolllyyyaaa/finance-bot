package inmem_config

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/util/file"
)

type SecretManager struct {
	data map[string]string
}

func NewSecretManager(readerModel file.Reader, filePath string) (*SecretManager, error) {
	mp, err := readerModel.ReadToMap(filePath)
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
