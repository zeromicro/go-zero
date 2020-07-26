package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"zero/core/discov"
	"zero/example/rpc/remote/unary"
	"zero/rpcx"
)

const timeFormat = "15:04:05"

func main() {
	flag.Parse()

	client := rpcx.MustNewClient(rpcx.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost:2379"},
			Key:   "rpcx",
		},
	})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			conn, ok := client.Next()
			if !ok {
				time.Sleep(time.Second)
				break
			}

			greet := unary.NewGreeterClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			resp, err := greet.Greet(ctx, &unary.Request{
				Name: "kevin",
			})
			if err != nil {
				fmt.Printf("%s X %s\n", time.Now().Format(timeFormat), err.Error())
			} else {
				fmt.Printf("%s => %s\n", time.Now().Format(timeFormat), resp.Greet)
			}
			cancel()
		}
	}
}
