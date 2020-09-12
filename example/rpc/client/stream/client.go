package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/example/rpc/remote/stream"
	"github.com/tal-tech/go-zero/rpcx"
)

const name = "kevin"

var key = flag.String("key", "rpcx", "the key on etcd")

func main() {
	flag.Parse()

	client, err := rpcx.NewClientNoAuth(discov.EtcdConf{
		Hosts: []string{"localhost:2379"},
		Key:   *key,
	})
	if err != nil {
		log.Fatal(err)
	}

	conn := client.Conn()
	greet := stream.NewStreamGreeterClient(conn)
	stm, err := greet.Greet(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			resp, err := stm.Recv()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("=>", resp.Greet)
		}
	}()

	for i := 0; i < 3; i++ {
		fmt.Println("<=", name)
		if err = stm.Send(&stream.StreamReq{
			Name: name,
		}); err != nil {
			log.Fatal(err)
		}
	}
}
