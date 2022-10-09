package config

import (
    "github.com/suyuan32/simple-admin-tools/plugins/registry/consul"
    "github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
}

