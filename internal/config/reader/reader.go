package reader

type Reader interface {
	Read(filepath string) (map[interface{}]interface{}, error)

	ToStringMap(raw map[interface{}]interface{}) (map[string]string, error)
}
