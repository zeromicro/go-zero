package main

import (
	"fmt"

	"zero/core/stores/redis"
	"zero/dq"
)

func main() {
	consumer := dq.NewConsumer(dq.DqConf{
		Beanstalks: []dq.Beanstalk{
			{
				Endpoint: "localhost:11300",
				Tube:     "tube",
			},
			{
				Endpoint: "localhost:11301",
				Tube:     "tube",
			},
			{
				Endpoint: "localhost:11302",
				Tube:     "tube",
			},
			{
				Endpoint: "localhost:11303",
				Tube:     "tube",
			},
			{
				Endpoint: "localhost:11304",
				Tube:     "tube",
			},
		},
		Redis: redis.RedisConf{
			Host: "localhost:6379",
			Type: redis.NodeType,
		},
	})
	consumer.Consume(func(body []byte) {
		fmt.Println(string(body))
	})
}
