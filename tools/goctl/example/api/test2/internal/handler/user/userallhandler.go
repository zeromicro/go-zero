package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test2/internal/logic/user"
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test2/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test2/internal/types"
)

func UserAllHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := user.NewUserAllLogic(r.Context(), svcCtx)
		resp, err := l.UserAll(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
