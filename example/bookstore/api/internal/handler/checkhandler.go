package handler

import (
	"bookstore/api/internal/logic"
	"bookstore/api/internal/svc"
	"bookstore/api/internal/types"
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func checkHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewCheckLogic(r.Context(), ctx)
		resp, err := l.Check(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.WriteJson(w, http.StatusOK, resp)
		}
	}
}
