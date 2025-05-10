package fileserver

import (
	"net/http"
	pathpkg "path"
	"strings"
	"sync"
)

// Middleware returns a middleware that serves files from the given file system.
func Middleware(path string, fs http.FileSystem) func(http.HandlerFunc) http.HandlerFunc {
	fileServer := http.FileServer(fs)
	pathWithoutTrailSlash := ensureNoTrailingSlash(path)
	canServe := createServeChecker(path, fs)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if canServe(r) {
				r.URL.Path = r.URL.Path[len(pathWithoutTrailSlash):]
				fileServer.ServeHTTP(w, r)
			} else {
				next(w, r)
			}
		}
	}
}

func createFileChecker(fs http.FileSystem) func(string) bool {
	var lock sync.RWMutex
	fileChecker := make(map[string]bool)

	return func(path string) bool {
		// Emulate http.Dir.Open’s path normalization for embed.FS.Open.
		// http.FileServer redirects any request ending in "/index.html"
		// to the same path without the final "index.html".
		// So the path here may be empty or end with a "/".
		// http.Dir.Open uses this logic to clean the path,
		// correctly handling those two cases.
		// embed.FS doesn’t perform this normalization, so we apply the same logic here.
		path = pathpkg.Clean("/" + path)[1:]
		if path == "" {
			path = "."
		}

		lock.RLock()
		exist, ok := fileChecker[path]
		lock.RUnlock()
		if ok {
			return exist
		}

		lock.Lock()
		defer lock.Unlock()

		file, err := fs.Open(path)
		exist = err == nil
		fileChecker[path] = exist
		if err != nil {
			return false
		}

		_ = file.Close()
		return true
	}
}

func createServeChecker(path string, fs http.FileSystem) func(r *http.Request) bool {
	pathWithTrailSlash := ensureTrailingSlash(path)
	fileChecker := createFileChecker(fs)

	return func(r *http.Request) bool {
		return r.Method == http.MethodGet &&
			strings.HasPrefix(r.URL.Path, pathWithTrailSlash) &&
			fileChecker(r.URL.Path[len(pathWithTrailSlash):])
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
