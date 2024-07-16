package handler

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

func fileServerHandler(subFS fs.FS, fsHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { // default path
		filePath := path.Clean(r.URL.Path)
		if filePath == rootPath {
			filePath = indexHtml
		} else {
			filePath = strings.TrimPrefix(filePath, slash)
		}

		file, err := subFS.Open(filePath)
		switch {
		case err == nil:
			fsHandler.ServeHTTP(w, r)
			_ = file.Close()
			return
		case os.IsNotExist(err):
			r.URL.Path = "/" // all virtual routes in react app means visit index.html
			fsHandler.ServeHTTP(w, r)
			return
		default:
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
	}
}
