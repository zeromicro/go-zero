package config

import (
	"zero/rest"
	"zero/rpcx"
)

type Config struct {
	rest.RestConf
	Rpc rpcx.RpcClientConf
}
