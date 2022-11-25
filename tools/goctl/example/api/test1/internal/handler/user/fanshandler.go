package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test1/internal/logic/user"
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test1/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test1/internal/types"
)

func FansHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FansRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := user.NewFansLogic(r.Context(), svcCtx)
		resp, err := l.Fans(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
