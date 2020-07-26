package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"zero/core/conf"
	"zero/example/rpc/remote/unary"
	"zero/rpcx"
)

var configFile = flag.String("f", "config.json", "the config file")

func main() {
	flag.Parse()

	var c rpcx.RpcClientConf
	conf.MustLoad(*configFile, &c)
	client := rpcx.MustNewClient(c)
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			conn, ok := client.Next()
			if !ok {
				log.Fatal("no server")
			}

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
