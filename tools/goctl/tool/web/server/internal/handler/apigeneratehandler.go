package handler

import (
	xhttp "github.com/zeromicro/x/http"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/logic"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/types"
)

func apiGenerateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.APIGenerateRequest
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := logic.NewApiGenerateLogic(r.Context(), svcCtx)
		resp, err := l.ApiGenerate(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
