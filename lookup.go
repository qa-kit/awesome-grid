package main

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

// Lookup provides addresses of pod in k8s cluster by sessiion
type Lookup struct {
	pool   *Pool
	config *Config
}

// SessionLength length of selenium session id
const SessionLength = 32

// Resolve finds ip of pod by session id, specified in request
func (l Lookup) Resolve(request *http.Request) (string, error) {
	log.Println("processing with existing session")
	path := request.URL.Path
	sessionID := strings.ReplaceAll(path, "/wd/hub/session/", "")
	if len(sessionID) < SessionLength {
		return "", errors.New("invalid session in path " + path)
	}

	sessionID = sessionID[:SessionLength]
	IP, exists := l.pool.IP(sessionID)
	if !exists {
		return "", errors.New("pod for session " + sessionID + " not found")
	}

	return "http://" + IP + ":" + l.config.PodPort, nil
}
