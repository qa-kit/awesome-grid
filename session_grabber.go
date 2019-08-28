package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"strings"
)

// SessionGrabber struct
type SessionGrabber struct {
	pool   *Pool
	config *Config
}

// RoundTrip processes reponses from web driver
func (g *SessionGrabber) grab(host string, bodyBytes []byte) error {
	webDriverResponse := CreateSessionResponse{}
	if err := json.Unmarshal(bodyBytes, &webDriverResponse); err != nil {
		return errors.New("parsing response, " + err.Error())
	}

	url, err := url.Parse("http://" + host)
	if err != nil {
		return errors.New("parsing " + host + ", " + err.Error())
	}
	// Removing port from host
	host = strings.ReplaceAll(url.Host, ":"+g.config.PodPort, "")
	log.Println("new session " + webDriverResponse.SessionID + " bind with " + host)
	g.pool.AddSession(webDriverResponse.SessionID, host)

	return nil
}
