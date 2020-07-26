package config

import (
	"zero/ngin"
	"zero/rpcx"
)

type Config struct {
	ngin.NgConf
	Rpc rpcx.RpcClientConf
}
