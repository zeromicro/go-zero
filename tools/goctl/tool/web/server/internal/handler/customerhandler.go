package handler

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
)

const (
	slash     = "/"
	rootPath  = "/"
	distDir   = "public"
	indexHtml = "index.html"
)

//go:embed public/*
var staticFiles embed.FS

func RegisterCustomHandlers(serverCtx *svc.ServiceContext) error {
	subFS, err := fs.Sub(staticFiles, distDir)
	if err != nil {
		return err
	}

	fileServer := http.FileServer(http.FS(subFS))
	http.HandleFunc("/", fileServerHandler(subFS, fileServer))
	http.HandleFunc("/api/generate", apiGenerateHandler(serverCtx))
	http.HandleFunc("/api/request/body/parse", requestBodyParseHandler(serverCtx))
	return nil
}
