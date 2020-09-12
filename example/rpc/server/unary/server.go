package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/example/rpc/remote/unary"
	"github.com/tal-tech/go-zero/rpcx"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/config.json", "the config file")

type GreetServer struct {
	lock     sync.Mutex
	alive    bool
	downTime time.Time
}

func NewGreetServer() *GreetServer {
	return &GreetServer{
		alive: true,
	}
}

func (gs *GreetServer) Greet(ctx context.Context, req *unary.Request) (*unary.Response, error) {
	fmt.Println("=>", req)

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &unary.Response{
		Greet: "hello from " + hostname,
	}, nil
}

func main() {
	flag.Parse()

	var c rpcx.RpcServerConf
	conf.MustLoad(*configFile, &c)

	server := rpcx.MustNewServer(c, func(grpcServer *grpc.Server) {
		unary.RegisterGreeterServer(grpcServer, NewGreetServer())
	})
	server.Start()
}
