package handler

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
)

const gzipEncoding = "gzip"

// GunzipHandler returns a middleware to gunzip http request body.
func GunzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get(httpx.ContentEncoding), gzipEncoding) {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			r.Body = reader
		}

		next.ServeHTTP(w, r)
	})
}
