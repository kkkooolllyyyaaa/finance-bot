package file

type Reader interface {
	Read(filepath string) (map[any]any, error)

	ToStringMap(raw map[any]any) (map[string]string, error)

	ReadToMap(filepath string) (map[string]string, error)
}
