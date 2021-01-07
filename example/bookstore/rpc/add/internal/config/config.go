package config

import (
	"github.com/3Rivers/go-zero/core/stores/cache"
	"github.com/3Rivers/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	Cache      cache.CacheConf
}
