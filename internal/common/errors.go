package common

import "github.com/pkg/errors"

var ErrIncorrectArgsCount = errors.New("Incorrect count of arguments")
var ErrIncorrectArgument = errors.New("Argument isn't fit format")
var ErrIncorrectUserID = errors.Wrap(ErrIncorrectArgument, "Can't parse userID")
