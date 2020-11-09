package rest

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func monitor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.Handler()
		h.ServeHTTP(w, r)
	}
}
