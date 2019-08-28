package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qa-kit/awesome-grid/cleaner"
	"github.com/qa-kit/awesome-grid/cluster"
	configPkg "github.com/qa-kit/awesome-grid/config"
	"github.com/qa-kit/awesome-grid/creator"
	"github.com/qa-kit/awesome-grid/lookup"
	poolPkg "github.com/qa-kit/awesome-grid/pool"
	"github.com/qa-kit/awesome-grid/session"
)

const (
	// ConfigPath path to config
	ConfigPath = "config.json"
)

func main() {
	var (
		r       = mux.NewRouter()
		config  = &configPkg.Config{}
		pool    = &poolPkg.Pool{}
		cluster = &cluster.Kubernetes{
			Config: config,
		}
		sessionGrabber    = session.New(pool, config)
		newSessionHandler = &ProxyHandler{
			resolver: creator.New(config, cluster, pool, &cleaner.Cleaner{}),
			transport: &Transport{
				callback:  sessionGrabber.Grab,
				roundTrip: http.DefaultTransport.RoundTrip,
			},
		}
		existSessionHandler = &ProxyHandler{
			resolver: lookup.New(pool, config),
		}
	)

	// New session handler
	r.HandleFunc("/wd/hub/session", newSessionHandler.Handle).Methods("POST")

	// Existing session handlerƒ
	r.PathPrefix("/").Handler(http.HandlerFunc(existSessionHandler.Handle))
	http.Handle("/", r)

	// Reading configs
	if err := config.Read(ConfigPath); err != nil {
		log.Fatal(err.Error())
	}

	// Creating clients
	if err := cluster.CreateClient(); err != nil {
		log.Fatal(err.Error())
	}

	if err := http.ListenAndServe(config.Listen, nil); err != nil {
		log.Fatal(err.Error())
	}
}
