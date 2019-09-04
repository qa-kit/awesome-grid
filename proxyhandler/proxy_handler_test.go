package proxyhandler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qa-kit/awesome-grid/transport"
)

type FakeResolver struct{}

func (c FakeResolver) Resolve(request *http.Request) (string, error) {
	return "", nil
}

func TestHandle(t *testing.T) {
	transport := transport.New(
		func(host string, bodyBytes []byte) error {
			return nil
		},
		func(request *http.Request) (*http.Response, error) {
			response := &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString("")),
			}
			return response, nil
		},
	)

	handler := New(FakeResolver{}, transport)

	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("creating http request error %s", err.Error())
		return
	}
	w := httptest.NewRecorder()
	handler.Handle(w, request)
}
