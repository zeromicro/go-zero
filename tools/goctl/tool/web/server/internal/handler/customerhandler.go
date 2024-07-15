package handler

import (
	"embed"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

//go:embed public/*
var staticFiles embed.FS

func RegisterCustomHandlers(serverCtx *svc.ServiceContext) {
	subFS, err := fs.Sub(staticFiles, "public")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(subFS))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // default path
		filePath := path.Clean(r.URL.Path)
		if filePath == "/" {
			filePath = "index.html"
		} else {
			filePath = strings.TrimPrefix(filePath, "/")
		}
		file, err := subFS.Open(filePath)
		switch {
		case err == nil:
			fileServer.ServeHTTP(w, r)
			file.Close()
			return
		case os.IsNotExist(err):
			r.URL.Path = "/" // all virtual routes in react app means visit index.html
			fileServer.ServeHTTP(w, r)
			return
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/api/generate", apiGenerateHandler(serverCtx))
	http.HandleFunc("/api/request/body/parse", requestBodyParseHandler(serverCtx))
}

type IndexFS struct {
	file fs.File
	fs   fs.FS
}

func NewIndexFS(fs fs.FS) *IndexFS {
	return &IndexFS{
		fs: fs,
	}
}
func (i *IndexFS) Open(name string) (http.File, error) {
	indexFS := &fileFS{i.fs}
	return http.FS(indexFS).Open("index.html")
}

type fileFS struct {
	fs fs.FS
}

func (f *fileFS) Open(name string) (fs.File, error) {
	file, err := f.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}
