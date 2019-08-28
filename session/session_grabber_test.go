package session

import (
	"testing"

	"github.com/qa-kit/awesome-grid/config"
	poolPkg "github.com/qa-kit/awesome-grid/pool"
)

func TestGrabSuccess(t *testing.T) {
	var (
		host     = "127.0.0.1"
		response = []byte("{\"sessionId\":\"id\"}")
		pool     = &poolPkg.Pool{}
		config   = &config.Config{PodPort: "4444"}
	)
	grabber := SessionGrabber{
		pool:   pool,
		config: config,
	}
	err := grabber.Grab(host, response)
	if err != nil {
		t.Errorf("expected empty error instead %s", err.Error())
	}
	found, exists := pool.IP("id")
	if found != host {
		t.Errorf("expected found '%s' instead '%s'", found, host)
	}
	if !exists {
		t.Errorf("expected exists %t instead %t", exists, true)
	}
}

func TestGrabFailParseResponse(t *testing.T) {
	var (
		host         = "http://127.0.0.1"
		errorMessage = "parsing response, unexpected end of JSON input"
		response     = []byte("")
		pool         = &poolPkg.Pool{}
	)
	grabber := SessionGrabber{
		pool: pool,
	}
	err := grabber.Grab(host, response)
	if err == nil {
		t.Errorf("expected error")
		return
	}
	if err.Error() != errorMessage {
		t.Errorf("expected error '%s' instead '%s'", err.Error(), errorMessage)
	}
}
