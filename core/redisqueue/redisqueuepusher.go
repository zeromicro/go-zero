package redisqueue

import (
	"fmt"
	"time"

	"zero/core/jsonx"
	"zero/core/logx"
	"zero/core/queue"
	"zero/core/stores/redis"
)

type (
	PusherOption func(p queue.QueuePusher) queue.QueuePusher

	RedisQueuePusher struct {
		name  string
		store *redis.Redis
		key   string
	}
)

func NewPusher(store *redis.Redis, key string, opts ...PusherOption) queue.QueuePusher {
	var pusher queue.QueuePusher = &RedisQueuePusher{
		name:  fmt.Sprintf("%s/%s/%s", store.Type, store.Addr, key),
		store: store,
		key:   key,
	}

	for _, opt := range opts {
		pusher = opt(pusher)
	}

	return pusher
}

func (saver *RedisQueuePusher) Name() string {
	return saver.name
}

func (saver *RedisQueuePusher) Push(message string) error {
	_, err := saver.store.Rpush(saver.key, message)
	if nil != err {
		return err
	}

	logx.Infof("<= %s", message)
	return nil
}

func WithTime() PusherOption {
	return func(p queue.QueuePusher) queue.QueuePusher {
		return timedQueuePusher{
			pusher: p,
		}
	}
}

type timedQueuePusher struct {
	pusher queue.QueuePusher
}

func (p timedQueuePusher) Name() string {
	return p.pusher.Name()
}

func (p timedQueuePusher) Push(message string) error {
	tm := TimedMessage{
		Time:    time.Now().Unix(),
		Payload: message,
	}

	if content, err := jsonx.Marshal(tm); err != nil {
		return err
	} else {
		return p.pusher.Push(string(content))
	}
}
