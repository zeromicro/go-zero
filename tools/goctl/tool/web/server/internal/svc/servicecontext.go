package svc

import (
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/config"
)

type ServiceContext struct {
	c *config.Config
}

func NewServiceContext(c *config.Config) *ServiceContext {
	return &ServiceContext{
		c: c,
	}
}
