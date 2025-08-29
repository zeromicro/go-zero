package svc

import (
    {{.importPackages}}
)

type ServiceContext struct {
	Config {{.config}}
}

func NewServiceContext(c {{.config}}) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
