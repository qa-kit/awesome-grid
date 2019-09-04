package proxyhandler

import "net/http"

// Resolver interface
type Resolver interface {
	Resolve(request *http.Request) (string, error)
}
