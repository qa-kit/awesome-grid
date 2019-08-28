package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FakeResolver struct{}

func (c FakeResolver) Resolve(request *http.Request) (string, error) {
	return "", nil
}

func TestHandle(t *testing.T) {
	transport := Transport{
		callback: func(host string, bodyBytes []byte) error {
			return nil
		},
		roundTrip: func(request *http.Request) (*http.Response, error) {
			response := &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}
			return response, nil
		},
	}
	handler := ProxyHandler{
		resolver:  FakeResolver{},
		transport: &transport,
	}
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("creating http request error %s", err.Error())
		return
	}
	w := httptest.NewRecorder()
	handler.Handle(w, request)
}
