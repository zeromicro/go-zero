package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"zero/core/cmdline"
	"zero/kq"
)

type message struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Payload string `json:"message"`
}

func main() {
	pusher := kq.NewPusher([]string{
		"172.16.56.64:19092",
		"172.16.56.65:19092",
		"172.16.56.66:19092",
	}, "kevin")

	ticker := time.NewTicker(time.Millisecond)
	for round := 0; round < 3; round++ {
		select {
		case <-ticker.C:
			count := rand.Intn(100)
			m := message{
				Key:     strconv.FormatInt(time.Now().UnixNano(), 10),
				Value:   fmt.Sprintf("%d,%d", round, count),
				Payload: fmt.Sprintf("%d,%d", round, count),
			}
			body, err := json.Marshal(m)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(body))
			if err := pusher.Push(string(body)); err != nil {
				log.Fatal(err)
			}
		}
	}

	cmdline.EnterToContinue()
}
