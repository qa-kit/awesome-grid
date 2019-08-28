package main

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ProxyHandler is flexible proxy handler
type ProxyHandler struct {
	resolver  Resolver
	transport *Transport
}

// Handle process request url
func (h *ProxyHandler) Handle(res http.ResponseWriter, request *http.Request) {
	log.Println("processing incoming request " + request.URL.Path)
	path, err := h.resolver.Resolve(request)
	if err != nil {
		h.HandleError(res, err)
		return
	}

	log.Println("proxing " + request.URL.Path + " to " + path)
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

	log.Println("proxing detail: " + request.Method + " " + request.URL.Scheme + " " + request.Host + request.URL.Path)
	proxy.ServeHTTP(res, request)
}

// HandleError writes error response
func (h *ProxyHandler) HandleError(res http.ResponseWriter, err error) {
	http.Error(res, err.Error(), 500)
}
