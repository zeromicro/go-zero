package cache

import "github.com/zeromicro/go-zero/core/stores/redis"

type (
	// A ClusterConf is the config of a redis cluster that used as cache.
	ClusterConf []NodeConf

	// A NodeConf is the config of a redis node that used as cache.
	NodeConf struct {
		redis.RedisConf `mapstructure:",squash"`
		Weight int `json:",default=100"`
	}
)
