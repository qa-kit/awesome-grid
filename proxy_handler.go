package main

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	logger "github.com/sirupsen/logrus"
)

// ProxyHandler is flexible proxy handler
type ProxyHandler struct {
	resolver  Resolver
	transport *Transport
}

// Handle process request url
func (h *ProxyHandler) Handle(res http.ResponseWriter, request *http.Request) {
	logger.Infof("processing incoming request %s", request.URL.Path)
	path, err := h.resolver.Resolve(request)
	if err != nil {
		h.HandleError(res, err)
		return
	}

	logger.Infof("proxing %s to %s", request.URL.Path, path)
	var target *url.URL
	target, err = url.Parse(path)
	if err != nil {
		h.HandleError(res, errors.New("proxy response, "+err.Error()))
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	if h.transport != nil {
		proxy.Transport = h.transport
	}
	request.URL.Host = target.Host
	request.URL.Scheme = target.Scheme
	request.Host = target.Host

	logger.Infof("proxing detail: %s %s %s%s", request.Method, request.URL.Scheme, request.Host, request.URL.Path)
	proxy.ServeHTTP(res, request)
}

// HandleError writes error response
func (h *ProxyHandler) HandleError(res http.ResponseWriter, err error) {
	http.Error(res, err.Error(), 500)
}
