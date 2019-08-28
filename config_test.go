package main

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	config := Config{}
	err := config.Read("config/config.json")
	if err != nil {
		t.Errorf("expected empty error instead '%s'", err.Error())
	}
}
