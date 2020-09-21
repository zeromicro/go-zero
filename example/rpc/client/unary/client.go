package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/example/rpc/remote/unary"
	"github.com/tal-tech/go-zero/zrpc"
)

var configFile = flag.String("f", "config.json", "the config file")

func main() {
	flag.Parse()

	var c zrpc.RpcClientConf
	conf.MustLoad(*configFile, &c)
	client := zrpc.MustNewClient(c)
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			conn := client.Conn()
			greet := unary.NewGreeterClient(conn)
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
