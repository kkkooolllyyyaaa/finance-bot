package currency

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const url string = "https://cdn.cur.su/api/cbr.json"

var RequestError = errors.New("Can't perform request")
var NotOkResponseStatusError = errors.New("Got non 200 status code")

type Api struct{}

func (ca *Api) Request() (responseBody io.ReadCloser, err error) {
	response, err := http.Get(url)
	if err != nil {
		return responseBody, errors.Wrap(RequestError, "Error while performing get request to: "+url)
	}
	defer func() { _ = response.Close }()

	if response.StatusCode != 200 {
		return responseBody, errors.Wrap(NotOkResponseStatusError, "Got http status code: "+response.Status)
	}

	return response.Body, nil
}
