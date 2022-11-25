package svc

import (
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test1/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
