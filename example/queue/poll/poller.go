package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"zero/core/discov"
	"zero/core/lang"
	"zero/core/logx"
	"zero/core/service"
	"zero/core/stores/redis"
	"zero/rq"
)

var (
	redisHost  = flag.String("redis", "localhost:6379", "")
	redisType  = flag.String("type", "node", "")
	redisKey   = flag.String("key", "queue", "")
	producers  = flag.Int("producers", 1, "")
	dropBefore = flag.Int64("drop", 0, "messages before seconds to drop")
)

type Consumer struct {
	lock      sync.Mutex
	resources map[string]interface{}
}

func NewConsumer() *Consumer {
	return &Consumer{
		resources: make(map[string]interface{}),
	}
}

func (c *Consumer) Consume(msg string) error {
	fmt.Println("=>", msg)
	c.lock.Lock()
	defer c.lock.Unlock()

	c.resources[msg] = lang.Placeholder

	return nil
}

func (c *Consumer) OnEvent(event interface{}) {
	fmt.Printf("event: %+v\n", event)
}

func main() {
	flag.Parse()

	consumer := NewConsumer()
	q, err := rq.NewMessageQueue(rq.RmqConf{
		ServiceConf: service.ServiceConf{
			Name: "queue",
			Log: logx.LogConf{
				Path:     "logs",
				KeepDays: 3,
				Compress: true,
			},
		},
		Redis: redis.RedisKeyConf{
			RedisConf: redis.RedisConf{
				Host: *redisHost,
				Type: *redisType,
			},
			Key: *redisKey,
		},
		Etcd: discov.EtcdConf{
			Hosts: []string{
				"localhost:2379",
			},
			Key: "queue",
		},
		DropBefore:   *dropBefore,
		NumProducers: *producers,
	}, rq.WithHandler(consumer), rq.WithRenewId(time.Now().UnixNano()))
	if err != nil {
		log.Fatal(err)
	}
	defer q.Stop()

	q.Start()
}
