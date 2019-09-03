package transport

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	var (
		response      = "response"
		catchResponse = ""
	)
	transport := Transport{
		callback: func(host string, bodyBytes []byte) error {
			catchResponse = string(bodyBytes)
			return nil
		},
		roundTrip: func(request *http.Request) (*http.Response, error) {
			response := &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(response)),
			}
			return response, nil
		},
	}
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("creating http request error %s", err.Error())
	}

	transport.RoundTrip(request)
	if response != catchResponse {
		t.Errorf("expected response '%s' instead '%s'", response, catchResponse)
	}
}

func TestRoundTripFailProxy(t *testing.T) {
	var errorMessage = "proxy error"
	var expectedErrorMessage = "round trip, " + errorMessage
	transport := Transport{
		callback: func(host string, bodyBytes []byte) error {
			return nil
		},
		roundTrip: func(request *http.Request) (*http.Response, error) {
			return nil, errors.New(errorMessage)
		},
	}
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("creating http request error %s", err.Error())
	}

	_, err = transport.RoundTrip(request)
	if err == nil {
		t.Error("expected error, nil instead")
	}

	if err.Error() != expectedErrorMessage {
		t.Errorf("expected error message '%s' instead '%s'", err.Error(), expectedErrorMessage)
	}
}

func TestRoundTripFailCallback(t *testing.T) {
	var errorMessage = "callback error"
	var expectedErrorMessage = "callback response, " + errorMessage
	transport := Transport{
		callback: func(host string, bodyBytes []byte) error {
			return errors.New(errorMessage)
		},
		roundTrip: func(request *http.Request) (*http.Response, error) {
			response := &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}
			return response, nil
		},
	}
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("creating http request error %s", err.Error())
	}

	_, err = transport.RoundTrip(request)
	if err == nil {
		t.Error("expected error, nil instead")
	}

	if err.Error() != expectedErrorMessage {
		t.Errorf("expected error message '%s' instead '%s'", err.Error(), expectedErrorMessage)
	}
}
