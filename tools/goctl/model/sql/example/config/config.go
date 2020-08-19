package config

import (
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

type (
	Config struct {
		logx.LogConf
		Mysql struct {
			DataSource string
			Table      struct {
				User   string
				Course string
			}
		}
		CacheRedis cache.CacheConf
		Redis      redis.RedisConf
	}
)
