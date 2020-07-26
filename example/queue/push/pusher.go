package main

import (
	"log"
	"strconv"
	"time"

	"zero/core/discov"
	"zero/rq"

	"github.com/google/gops/agent"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}

	pusher, err := rq.NewPusher([]string{"localhost:2379"}, "queue", rq.WithConsistentStrategy(
		func(msg string) (string, string, error) {
			return msg, msg, nil
		}, discov.BalanceWithId()), rq.WithServerSensitive())
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; ; i++ {
		pusher.Push(strconv.Itoa(i))
		time.Sleep(time.Second)
	}
}
