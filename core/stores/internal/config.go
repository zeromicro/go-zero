package internal

import "zero/core/stores/redis"

type (
	ClusterConf []NodeConf

	NodeConf struct {
		redis.RedisConf
		Weight int `json:",default=100"`
	}
)
