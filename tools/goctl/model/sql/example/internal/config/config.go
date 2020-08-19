package config

import (
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/cache"
)

type (
	Config struct {
		logx.LogConf
		Mysql struct {
			DataSource string
			Table      struct {
				User string
			}
		}
		CacheRedis cache.CacheConf
	}
)
