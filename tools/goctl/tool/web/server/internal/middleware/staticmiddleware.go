package middleware

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

type StaticMiddleware struct {
	assets embed.FS
}

func NewStaticMiddleware(assets embed.FS) *StaticMiddleware {
	return &StaticMiddleware{}
}

func (m *StaticMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	fs, err := fs.Sub(m.assets, "static")
	if err != nil {
		log.Fatal(err)
	}

	path := "/static"
	fileServer := http.FileServer(http.FS(fs))
	pathWithTrailSlash := ensureTrailingSlash(path)
	pathWithoutTrailSlash := ensureNoTrailingSlash(path)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, pathWithTrailSlash) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, pathWithoutTrailSlash)
			fileServer.ServeHTTP(w, r)
		} else {
			next(w, r)
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
