package fileserver

import (
	"net/http"
	"strings"
)

func Middleware(path, dir string) func(http.HandlerFunc) http.HandlerFunc {
	fileServer := http.FileServer(http.Dir(dir))
	pathWithTrailSlash := ensureTrailingSlash(path)
	pathWithoutTrailSlash := ensureNoTrailingSlash(path)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, pathWithTrailSlash) {
				r.URL.Path = strings.TrimPrefix(r.URL.Path, pathWithoutTrailSlash)
				fileServer.ServeHTTP(w, r)
			} else {
				next(w, r)
			}
		}
	}
}

func ensureTrailingSlash(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	}

	return path + "/"
}

func ensureNoTrailingSlash(path string) string {
	if strings.HasSuffix(path, "/") {
		return path[:len(path)-1]
	}

	return path
}
