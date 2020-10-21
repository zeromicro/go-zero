package rest

import (
	"net/http"
	"strings"
)

const (
	allowOrigin  = "Access-Control-Allow-Origin"
	allOrigin    = "*"
	allowMethods = "Access-Control-Allow-Methods"
	allowHeaders = "Access-Control-Allow-Headers"
	headers      = "Content-Type, Content-Length, Origin"
	methods      = "GET, HEAD, POST, PATCH, PUT, DELETE"
	separator    = ", "
)

func CorsHandler(origins ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(origins) > 0 {
			w.Header().Set(allowOrigin, strings.Join(origins, separator))
		} else {
			w.Header().Set(allowOrigin, allOrigin)
		}
		w.Header().Set(allowMethods, methods)
		w.Header().Set(allowHeaders, headers)
		w.WriteHeader(http.StatusNoContent)
	})
}
