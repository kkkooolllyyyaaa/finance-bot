package reader

import "github.com/pkg/errors"

type Reader interface {
	Read(filepath string) (map[any]any, error)

	ToStringMap(raw map[any]any) (map[string]string, error)
}

type Model struct {
	reader Reader
}

func New(reader Reader) *Model {
	return &Model{reader: reader}
}

func (m *Model) ReadToMap(filepath string) (map[string]string, error) {
	raw, err := m.reader.Read(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "reader Read")
	}

	data, err := m.reader.ToStringMap(raw)
	if err != nil {
		return nil, errors.Wrap(err, "reader ToStringMap")
	}

	return data, nil
}
