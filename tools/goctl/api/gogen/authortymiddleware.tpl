package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/suyuan32/simple-admin-core/pkg/i18n"
	"github.com/suyuan32/simple-admin-core/pkg/msg/logmsg"

	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type AuthorityMiddleware struct {
	Cbn *casbin.Enforcer
	Rds *redis.Redis
}

func NewAuthorityMiddleware(cbn *casbin.Enforcer, rds *redis.Redis) *AuthorityMiddleware {
	return &AuthorityMiddleware{
		Cbn: cbn,
		Rds: rds,
	}
}

func (m *AuthorityMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the path
		obj := r.URL.Path
		// get the method
		act := r.Method
		// get the role id
		roleId := r.Context().Value("roleId").(json.Number).String()
		// check the role status
		roleStatus, err := m.Rds.Hget("roleData", fmt.Sprintf("%s_status", roleId))
		if err != nil {
			logx.Errorw(logmsg.RedisError, logx.Field("detail", err.Error()))
			httpx.Error(w, errorx.NewApiErrorWithoutMsg(http.StatusUnauthorized))
			return
		} else if roleStatus == "0" {
			logx.Errorw("role is on forbidden status", logx.Field("roleId", roleId))
			httpx.Error(w, errorx.NewApiError(http.StatusBadRequest, "role.roleForbidden"))
			return
		}

		// check jwt blacklist
		res, err := m.Rds.Get("token_" + r.Header.Get("Authorization"))
		if err != nil {
			logx.Errorw("redis error in jwt", logx.Field("detail", err.Error()))
			httpx.Error(w, errorx.NewApiError(http.StatusInternalServerError, err.Error()))
			return
		}
		if res == "1" {
			logx.Errorw("token in blacklist", logx.Field("detail", r.Header.Get("Authorization")))
			httpx.Error(w, errorx.NewApiErrorWithoutMsg(http.StatusUnauthorized))
			return
		}

		sub := roleId
		result, err := m.Cbn.Enforce(sub, obj, act)
		if err != nil {
			logx.Errorw("casbin enforce error", logx.Field("detail", err.Error()))
			httpx.Error(w, errorx.NewApiError(http.StatusInternalServerError, i18n.ApiRequestFailed))
			return
		}
		if result {
			logx.Infow("HTTP/HTTPS Request", logx.Field("UUID", r.Context().Value("userId").(string)),
				logx.Field("path", obj), logx.Field("method", act))
			next(w, r)
			return
		} else {
			logx.Errorw("the role is not permitted to access the API", logx.Field("roleId", roleId),
				logx.Field("path", obj), logx.Field("method", act))
			httpx.Error(w, errorx.NewApiErrorWithoutMsg(http.StatusForbidden))
			return
		}
	}
}
