package redisqueue

import (
	"fmt"
	"sync"
	"time"

	"zero/core/jsonx"
	"zero/core/logx"
	"zero/core/queue"
	"zero/core/stores/redis"
)

const (
	logIntervalMillis  = 1000
	retryRedisInterval = time.Second
)

type (
	ProducerOption func(p queue.Producer) queue.Producer

	RedisQueueProducer struct {
		name      string
		store     *redis.Redis
		key       string
		redisNode redis.ClosableNode
		listeners []queue.ProduceListener
	}
)

func NewProducerFactory(store *redis.Redis, key string, opts ...ProducerOption) queue.ProducerFactory {
	return func() (queue.Producer, error) {
		return newProducer(store, key, opts...)
	}
}

func (p *RedisQueueProducer) AddListener(listener queue.ProduceListener) {
	p.listeners = append(p.listeners, listener)
}

func (p *RedisQueueProducer) Name() string {
	return p.name
}

func (p *RedisQueueProducer) Produce() (string, bool) {
	lessLogger := logx.NewLessLogger(logIntervalMillis)

	for {
		value, ok, err := p.store.BlpopEx(p.redisNode, p.key)
		if err == nil {
			return value, ok
		} else if err == redis.Nil {
			// timed out without elements popped
			continue
		} else {
			lessLogger.Errorf("Error on blpop: %v", err)
			p.waitForRedisAvailable()
		}
	}
}

func newProducer(store *redis.Redis, key string, opts ...ProducerOption) (queue.Producer, error) {
	redisNode, err := redis.CreateBlockingNode(store)
	if err != nil {
		return nil, err
	}

	var producer queue.Producer = &RedisQueueProducer{
		name:      fmt.Sprintf("%s/%s/%s", store.Type, store.Addr, key),
		store:     store,
		key:       key,
		redisNode: redisNode,
	}

	for _, opt := range opts {
		producer = opt(producer)
	}

	return producer, nil
}

func (p *RedisQueueProducer) resetRedisConnection() error {
	if p.redisNode != nil {
		p.redisNode.Close()
		p.redisNode = nil
	}

	redisNode, err := redis.CreateBlockingNode(p.store)
	if err != nil {
		return err
	}

	p.redisNode = redisNode
	return nil
}

func (p *RedisQueueProducer) waitForRedisAvailable() {
	var paused bool
	var pauseOnce sync.Once

	for {
		if err := p.resetRedisConnection(); err != nil {
			pauseOnce.Do(func() {
				paused = true
				for _, listener := range p.listeners {
					listener.OnProducerPause()
				}
			})
			logx.Errorf("Error occurred while connect to redis: %v", err)
			time.Sleep(retryRedisInterval)
		} else {
			break
		}
	}

	if paused {
		for _, listener := range p.listeners {
			listener.OnProducerResume()
		}
	}
}

func TimeSensitive(seconds int64) ProducerOption {
	return func(p queue.Producer) queue.Producer {
		if seconds > 0 {
			return autoDropQueueProducer{
				seconds:  seconds,
				producer: p,
			}
		}

		return p
	}
}

type autoDropQueueProducer struct {
	seconds  int64 // seconds before to drop
	producer queue.Producer
}

func (p autoDropQueueProducer) AddListener(listener queue.ProduceListener) {
	p.producer.AddListener(listener)
}

func (p autoDropQueueProducer) Produce() (string, bool) {
	lessLogger := logx.NewLessLogger(logIntervalMillis)

	for {
		content, ok := p.producer.Produce()
		if !ok {
			return "", false
		}

		var timedMsg TimedMessage
		if err := jsonx.UnmarshalFromString(content, &timedMsg); err != nil {
			lessLogger.Errorf("invalid timedMessage: %s, error: %s", content, err.Error())
			continue
		}

		if timedMsg.Time+p.seconds < time.Now().Unix() {
			lessLogger.Errorf("expired timedMessage: %s", content)
		}

		return timedMsg.Payload, true
	}
}
