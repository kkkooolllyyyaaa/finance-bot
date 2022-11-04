package yaml

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type YamlReader struct{}

func New() *YamlReader {
	return &YamlReader{}
}

func (r *YamlReader) Read(filepath string) (result map[any]any, err error) {
	rawYAML, err := os.ReadFile(filepath)
	if err != nil {
		return result, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &result)
	if err != nil {
		return result, errors.Wrap(err, "parsing yaml")
	}

	return result, nil
}

func (r *YamlReader) ToStringMap(raw map[any]any) (map[string]string, error) {
	result := make(map[string]string, len(raw))

	for k, v := range raw {
		kstr, ok := k.(string)
		if !ok {
			return nil, errors.New("Can't convert key to string")
		}
		vstr, ok := v.(string)
		if !ok {
			return nil, errors.New("Can't convert value to string")
		}

		result[kstr] = vstr
	}
	return result, nil
}

func (r *YamlReader) ReadToMap(filepath string) (map[string]string, error) {
	raw, err := r.Read(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "reader Read")
	}

	data, err := r.ToStringMap(raw)
	if err != nil {
		return nil, errors.Wrap(err, "reader ToStringMap")
	}

	return data, nil
}
