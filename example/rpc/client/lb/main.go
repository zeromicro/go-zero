package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/example/rpc/remote/unary"
	"github.com/tal-tech/go-zero/rpcx"
)

var lb = flag.String("t", "direct", "the load balancer type")

func main() {
	flag.Parse()

	var cli rpcx.Client
	switch *lb {
	case "direct":
		cli = rpcx.MustNewClient(rpcx.RpcClientConf{
			Endpoints: []string{
				"localhost:3456",
				"localhost:3457",
			},
		})
	case "discov":
		cli = rpcx.MustNewClient(rpcx.RpcClientConf{
			Etcd: discov.EtcdConf{
				Hosts: []string{"localhost:2379"},
				Key:   "rpcx",
			},
		})
	default:
		log.Fatal("bad load balancing type")
	}

	greet := unary.NewGreeterClient(cli.Conn())
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			resp, err := greet.Greet(context.Background(), &unary.Request{
				Name: "kevin",
			})
			if err != nil {
				fmt.Println("X", err.Error())
			} else {
				fmt.Println("=>", resp.Greet)
			}
		}
	}
}
