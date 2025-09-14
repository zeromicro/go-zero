// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1-alpha

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"workspace/internal/logic"
	"workspace/internal/svc"
	"workspace/internal/types"
)

func GreetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGreetLogic(r.Context(), svcCtx)
		resp, err := l.Greet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
