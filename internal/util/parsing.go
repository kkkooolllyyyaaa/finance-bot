package util

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var ErrNotACommand = errors.New("It's not a command")

func ParseCommand(text string) (command string, args []string, err error) {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	if !strings.HasPrefix(text, "/") {
		err = ErrNotACommand
		return
	}

	fields := strings.Fields(text)

	return fields[0], fields[1:], nil
}

func ParseInt64(number string) (int64, error) {
	return strconv.ParseInt(number, 10, 64)
}

func ParseInt(number string) (int, error) {
	numberInt64, err := ParseInt64(number)
	return int(numberInt64), err
}

func ParseFloat64(number string) (float64, error) {
	return strconv.ParseFloat(number, 64)
}
