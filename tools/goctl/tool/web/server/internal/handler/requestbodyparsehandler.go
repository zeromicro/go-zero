package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/logic"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/types"
	xhttp "github.com/zeromicro/x/http"
)

func requestBodyParseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ParseJsonRequest
		if err := httpx.Parse(r, &req); err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRequestBodyParseLogic(r.Context(), svcCtx)
		resp, err := l.RequestBodyParse(&req)
		if err != nil {
			xhttp.JsonBaseResponseCtx(r.Context(), w, err)
		} else {
			xhttp.JsonBaseResponseCtx(r.Context(), w, resp)
		}
	}
}
