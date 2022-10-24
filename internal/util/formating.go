package util

import "strconv"

type floatNumber interface {
	float32 | float64
}

func FormatFloat[T floatNumber](number T) string {
	return strconv.FormatFloat(float64(number), 'f', 2, 64)
}

type integerNumber interface {
	int | int8 | int16 | int32 | int64
}

func FormatInt[T integerNumber](number T) string {
	return strconv.FormatInt(int64(number), 10)
}
