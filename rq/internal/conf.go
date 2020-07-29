package internal

import (
	"zero/core/queue"
	"zero/core/stores/redis"
)

type RedisKeyConf struct {
	redis.RedisConf
	Key string `json:",optional"`
}

func (rkc RedisKeyConf) NewProducer(opts ...ProducerOption) (queue.Producer, error) {
	return newProducer(rkc.NewRedis(), rkc.Key, opts...)
}

func (rkc RedisKeyConf) NewPusher(opts ...PusherOption) queue.QueuePusher {
	return NewPusher(rkc.NewRedis(), rkc.Key, opts...)
}
