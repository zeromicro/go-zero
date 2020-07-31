package config

import (
	"zero/rest"
	"zero/rpcx"
)

type Config struct {
	rest.RtConf
	Rpc rpcx.RpcClientConf
}
