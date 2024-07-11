package handler

import (
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"net/http"
)

func RegisterCustomHandlers(serverCtx *svc.ServiceContext) {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("build"))))
	http.HandleFunc("/api/generate", apiGenerateHandler(serverCtx))
}
