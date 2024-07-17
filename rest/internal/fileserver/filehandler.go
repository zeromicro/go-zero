package fileserver

import (
	"net/http"
	"strings"
)

// Middleware returns a middleware that serves files from the given file system.
func Middleware(path string, fs http.FileSystem) func(http.HandlerFunc) http.HandlerFunc {
	fileServer := http.FileServer(fs)
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
