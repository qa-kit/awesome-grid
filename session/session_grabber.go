package session

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/qa-kit/awesome-grid/config"
	poolPkg "github.com/qa-kit/awesome-grid/pool"
	logger "github.com/sirupsen/logrus"
)

// SessionGrabber struct
type SessionGrabber struct {
	pool   *poolPkg.Pool
	config *config.Config
}

//New creates new SessionGrabber
func New(pool *poolPkg.Pool, config *config.Config) SessionGrabber {
	return SessionGrabber{pool: pool, config: config}
}

// Grab processes reponses from web driver
func (g *SessionGrabber) Grab(host string, bodyBytes []byte) error {
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
	logger.Infof("new session %s bind with %s", webDriverResponse.SessionID, host)
	g.pool.AddSession(webDriverResponse.SessionID, host)

	return nil
}
