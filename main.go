package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qa-kit/awesome-grid/cleaner"
	"github.com/qa-kit/awesome-grid/cluster"
	configPkg "github.com/qa-kit/awesome-grid/config"
	"github.com/qa-kit/awesome-grid/creator"
	"github.com/qa-kit/awesome-grid/lookup"
	poolPkg "github.com/qa-kit/awesome-grid/pool"
	"github.com/qa-kit/awesome-grid/proxyhandler"
	"github.com/qa-kit/awesome-grid/session"
	"github.com/qa-kit/awesome-grid/transport"
	logger "github.com/sirupsen/logrus"
)

const (
	// ConfigPath path to config
	ConfigPath = "config.json"
)

func main() {

	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})

	var (
		r       = mux.NewRouter()
		config  = &configPkg.Config{}
		pool    = &poolPkg.Pool{}
		cluster = &cluster.Kubernetes{
			Config: config,
		}
		sessionGrabber = session.New(pool, config)

		newSessionHandler = proxyhandler.New(
			creator.New(config, cluster, pool, &cleaner.Cleaner{}, http.Get),
			transport.New(sessionGrabber.Grab, http.DefaultTransport.RoundTrip),
		)

		existSessionHandler = proxyhandler.New(
			lookup.New(pool, config),
			nil,
		)
	)

	// New session handler
	r.HandleFunc("/wd/hub/session", newSessionHandler.Handle)
	// Healthcheck
	r.HandleFunc("/status/", func(res http.ResponseWriter, request *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte("{\"status\":\"ok\"}"))
	})
	// Existing session handler
	r.PathPrefix("/").Handler(http.HandlerFunc(existSessionHandler.Handle))

	http.Handle("/", r)

	// Reading configs
	if err := config.Read(ConfigPath); err != nil {
		logger.Fatal(err.Error())
	}

	// Creating clients
	if err := cluster.CreateClient(); err != nil {
		logger.Fatal(err.Error())
	}

	if err := http.ListenAndServe(config.Listen, nil); err != nil {
		logger.Fatal(err.Error())
	}
}
