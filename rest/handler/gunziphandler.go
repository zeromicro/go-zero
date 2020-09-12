package handler

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/tal-tech/go-zero/rest/httpx"
)

const gzipEncoding = "gzip"

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
