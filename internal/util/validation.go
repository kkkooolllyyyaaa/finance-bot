package util

type integer interface {
	int | int8 | int16 | int32 | int64
}

func ValidateNatural[T integer](number T) (ok bool) {
	return number > 0
}
