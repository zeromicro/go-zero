package svc

import (
	"{{.configPkg}}"{{if .callRPC}}
	"github.com/zeromicro/go-zero/zrpc"
	"{{.rpcClientPkg}}"{{end}}
)

type ServiceContext struct {
	Config   config.Config{{if .callRPC}}
	GreetRpc greet.Greet{{end}}
}

func NewServiceContext(c config.Config) *ServiceContext {
	{{if .callRPC}}client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "127.0.0.1:8080",
	}){{end}}
	return &ServiceContext{
		Config:   c,
		{{if .callRPC}}GreetRpc: greet.NewGreet(client),{{end}}
	}
}
