package reader

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type YamlReader struct{}

func New() *YamlReader {
	return &YamlReader{}
}

func (r *YamlReader) Read(filepath string) (result map[interface{}]interface{}, err error) {

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

func (r *YamlReader) ToStringMap(raw map[interface{}]interface{}) (map[string]string, error) {
	result := make(map[string]string, len(raw))

	for k, v := range raw {
		kstr := fmt.Sprintf("%v", k)
		vstr := fmt.Sprintf("%v", v)
		result[kstr] = vstr
	}
	return result, nil
}
