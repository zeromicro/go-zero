package rq

import (
	"zero/core/discov"
	"zero/core/service"
	"zero/core/stores/redis"
)

type RmqConf struct {
	service.ServiceConf
	Redis           redis.RedisKeyConf
	Etcd            discov.EtcdConf `json:",optional"`
	NumProducers    int             `json:",optional"`
	NumConsumers    int             `json:",optional"`
	Timeout         int64           `json:",optional"`
	DropBefore      int64           `json:",optional"`
	ServerSensitive bool            `json:",default=false"`
}
