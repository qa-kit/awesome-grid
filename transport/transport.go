package transport

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

// Transport struct
type Transport struct {
	callback  func(host string, bodyBytes []byte) error
	roundTrip func(request *http.Request) (*http.Response, error)
}

//New creates new transport
func New(
	callback func(host string, bodyBytes []byte) error,
	roundTrip func(request *http.Request) (*http.Response, error),
) *Transport {
	return &Transport{callback, roundTrip}
}

// RoundTrip processes reponses from web driver
func (t *Transport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := t.roundTrip(request)
	if err != nil {
		return nil, errors.New("round trip, " + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("reading response, " + err.Error())
	}

	if err = t.callback(request.Host, bodyBytes); err != nil {
		return nil, errors.New("callback response, " + err.Error())
	}

	response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return response, nil
}
