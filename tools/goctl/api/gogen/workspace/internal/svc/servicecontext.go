// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.1-alpha

package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"workspace/internal/config"
	"workspace/internal/middleware"
)

type ServiceContext struct {
	Config        config.Config
	TokenValidate rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		TokenValidate: middleware.NewTokenValidateMiddleware().Handle,
	}
}
