package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tal-tech/go-zero/core/discov"
)

var value = flag.String("v", "value", "the value")

func main() {
	flag.Parse()

	client := discov.NewPublisher([]string{"etcd.discovery:2379"}, "028F2C35852D", *value)
	if err := client.KeepAlive(); err != nil {
		log.Fatal(err)
	}
	defer client.Stop()

	for {
		time.Sleep(time.Second)
		fmt.Println(*value)
	}
}
