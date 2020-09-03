package config

import (
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/rpcx"
)

type Config struct {
	rest.RestConf
	Add   rpcx.RpcClientConf
	Check rpcx.RpcClientConf
}