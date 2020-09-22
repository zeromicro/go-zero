package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/example/rpc/remote/stream"
	"github.com/tal-tech/go-zero/zrpc"
)

const name = "kevin"

var key = flag.String("key", "zrpc", "the key on etcd")

func main() {
	flag.Parse()

	client, err := zrpc.NewClientNoAuth(discov.EtcdConf{
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

	var wg sync.WaitGroup
	go func() {
		for {
			resp, err := stm.Recv()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("=>", resp.Greet)
			wg.Done()
		}
	}()

	for i := 0; i < 3; i++ {
		wg.Add(1)
		fmt.Println("<=", name)
		if err = stm.Send(&stream.StreamReq{
			Name: name,
		}); err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()
}
