package handler

import (
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"io/fs"
	"log"
	"net/http"
)

func RegisterCustomHandlers(serverCtx *svc.ServiceContext) {
	fs, err := fs.Sub(serverCtx.Assets, "static")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(fs))))
	http.HandleFunc("/api/generate", apiGenerateHandler(serverCtx))
	http.HandleFunc("/api/request/body/parse", requestBodyParseHandler(serverCtx))
}
