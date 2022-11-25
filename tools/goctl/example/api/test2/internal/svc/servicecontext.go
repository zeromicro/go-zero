package svc

import (
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test2/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
