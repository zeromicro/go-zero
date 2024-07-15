package svc

import (
	"embed"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/config"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/middleware"
)

type ServiceContext struct {
	Config config.Config
	Static rest.Middleware
}

func NewServiceContext(c config.Config, assets embed.FS) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Static: middleware.NewStaticMiddleware(assets).Handle,
	}
}
