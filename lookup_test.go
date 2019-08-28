package main

import (
	"net/http"
	"testing"
)

func TestResolve(t *testing.T) {
	var (
		sessionID = "C20AD4D76FE97759AA27A0C99BFF6710"
		host      = "127.0.0.1"
		expected  = "http://127.0.0.1:4444"
		pool      = &Pool{}
		config    = &Config{PodPort: "4444"}
	)
	lookup := Lookup{
		pool:   pool,
		config: config,
	}
	pool.AddSession(sessionID, host)
	request, err := http.NewRequest("GET", "/wd/hub/session/"+sessionID, nil)
	target, err := lookup.Resolve(request)
	if err != nil {
		t.Errorf("expected empty error instead %s", err.Error())
	}
	if target != expected {
		t.Errorf("expected target '%s' instead '%s'", target, expected)
	}
}
