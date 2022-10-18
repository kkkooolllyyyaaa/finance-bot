package config

import (
	"github.com/pkg/errors"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/config/reader"
)

const secretsFile = "config/secrets/tg-bot-api.yaml"

type Secret struct {
	data map[string]string
}

func New(reader reader.Reader) (*Secret, error) {
	raw, err := reader.Read(secretsFile)
	if err != nil {
		return nil, errors.Wrap(err, "reader Read")
	}

	data, err := reader.ToStringMap(raw)
	if err != nil {
		return nil, errors.Wrap(err, "reader ToStringMap")
	}

	return &Secret{
		data: data,
	}, nil
}

func (s *Secret) Token() string {
	return s.data["token"]
}
