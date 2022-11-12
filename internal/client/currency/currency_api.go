package currency

import (
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const url string = "https://www.cbr-xml-daily.ru/latest.js"

var requestError = errors.New("Can't perform request")
var notOkResponseStatusError = errors.New("Got non 200 status code")

type Api struct {
	client http.Client
}

func NewCurrencyApi() *Api {
	return &Api{
		client: http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (ca *Api) Request() (responseBody io.ReadCloser, err error) {
	response, err := ca.client.Get(url)
	if err != nil {
		return responseBody, errors.Wrap(requestError, "Error while performing get request to: "+url)
	}
	defer func() { _ = response.Close }()

	if response.StatusCode != 200 {
		return responseBody, errors.Wrap(notOkResponseStatusError, "Got http status code: "+response.Status)
	}

	return response.Body, nil
}
