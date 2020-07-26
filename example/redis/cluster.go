package main

import (
	"flag"
	"log"

	"zero/core/logx"
	"zero/core/queue"
	"zero/core/service"
	"zero/core/stores/redis"
	"zero/rq"
)

var (
	host = flag.String("s", "10.24.232.63:7002", "server address")
	mode = flag.String("m", "queue", "cluster test mode")
)

type bridgeHandler struct {
	pusher queue.QueuePusher
}

func newBridgeHandler() rq.ConsumeHandler {
	return bridgeHandler{}
}

func (h bridgeHandler) Consume(str string) error {
	logx.Info("=>", str)
	return nil
}

func main() {
	flag.Parse()

	if *mode == "queue" {
		mq, err := rq.NewMessageQueue(rq.RmqConf{
			ServiceConf: service.ServiceConf{
				Log: logx.LogConf{
					Path: "logs",
				},
			},
			Redis: redis.RedisKeyConf{
				RedisConf: redis.RedisConf{
					Host: *host,
					Type: "cluster",
				},
				Key: "notexist",
			},
			NumProducers: 1,
		}, rq.WithHandler(newBridgeHandler()))
		if err != nil {
			log.Fatal(err)
		}
		defer mq.Stop()

		mq.Start()
	} else {
		rds := redis.NewRedis(*host, "cluster")
		rds.Llen("notexist")
		select {}
	}
}
