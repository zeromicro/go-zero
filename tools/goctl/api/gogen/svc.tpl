package svc

import (
	{{.configImport}}
	{{if .useI18n}}
	"github.com/suyuan32/simple-admin-core/pkg/i18n"{{end}}
    {{if .useCasbin}}
	"github.com/zeromicro/go-zero/core/stores/redis"
    "github.com/zeromicro/go-zero/rest"
    "github.com/casbin/casbin/v2"{{end}}
)

type ServiceContext struct {
	Config {{.config}}
	{{.middleware}}{{if .useCasbin}}Casbin    *casbin.Enforcer
	Authority rest.Middleware{{end}}
	{{if .useI18n}}Trans     *i18n.Translator{{end}}
}

func NewServiceContext(c {{.config}}) *ServiceContext {
{{if .useCasbin}}
    rds := redis.MustNewRedis(c.RedisConf)

    cbn := c.CasbinConf.MustNewCasbinWithRedisWatcher(c.DatabaseConf.Type, c.DatabaseConf.GetDSN(), c.RedisConf)
{{end}}
{{if .useI18n}}
    trans := &i18n.Translator{}
    trans.NewBundle(i18n2.LocaleFS)
    trans.NewTranslator()
{{end}}
	return &ServiceContext{
		Config: c,
		{{if .useCasbin}}Authority: middleware.NewAuthorityMiddleware(cbn, rds, trans).Handle,{{end}}
		{{if .useI18n}}Trans:     trans,{{end}}
	}
}
