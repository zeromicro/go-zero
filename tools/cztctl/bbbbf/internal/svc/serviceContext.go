package svc

import (
	"github.com/lerity-yao/go-zero/tools/cztctl/bbbbf/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
