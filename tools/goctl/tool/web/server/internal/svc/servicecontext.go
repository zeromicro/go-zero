package svc

import (
	"embed"
)

type ServiceContext struct {
	Assets embed.FS
}

func NewServiceContext(assets embed.FS) *ServiceContext {
	return &ServiceContext{
		Assets: assets,
	}
}
