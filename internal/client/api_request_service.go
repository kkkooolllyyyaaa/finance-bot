package client

import (
	"io"
)

type ApiRequestService interface {
	Request() (responseBody io.ReadCloser, err error)
}
