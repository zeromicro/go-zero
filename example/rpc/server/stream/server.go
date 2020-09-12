package main

import (
	"fmt"
	"io"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/example/rpc/remote/stream"
	"github.com/tal-tech/go-zero/rpcx"
	"google.golang.org/grpc"
)

type StreamGreetServer int

func (gs StreamGreetServer) Greet(s stream.StreamGreeter_GreetServer) error {
	ctx := s.Context()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("cancelled by client")
			return ctx.Err()
		default:
			req, err := s.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}

			fmt.Println("=>", req.Name)
			greet := "hello, " + req.Name
			fmt.Println("<=", greet)
			s.Send(&stream.StreamResp{
				Greet: greet,
			})
		}
	}
}

func main() {
	var c rpcx.RpcServerConf
	conf.MustLoad("etc/config.json", &c)

	server := rpcx.MustNewServer(c, func(grpcServer *grpc.Server) {
		stream.RegisterStreamGreeterServer(grpcServer, StreamGreetServer(0))
	})
	server.Start()
}
