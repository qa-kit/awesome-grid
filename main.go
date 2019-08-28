package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	// ConfigPath path to config
	ConfigPath = "config/config.json"
)

func main() {
	var (
		r       = mux.NewRouter()
		config  = &Config{}
		pool    = &Pool{}
		cluster = &Kubernetes{
			config: config,
		}
		sessionGrabber = SessionGrabber{
			pool:   pool,
			config: config,
		}
		newSessionHandler = &ProxyHandler{
			resolver: &Creator{
				config:  config,
				cluster: cluster,
				pool:    pool,
				cleaner: &Cleaner{},
			},
			transport: &Transport{
				callback:  sessionGrabber.grab,
				roundTrip: http.DefaultTransport.RoundTrip,
			},
		}
		existSessionHandler = &ProxyHandler{
			resolver: Lookup{
				pool:   pool,
				config: config,
			},
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
