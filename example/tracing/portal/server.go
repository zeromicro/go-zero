package main

import (
	"context"
	"flag"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/example/tracing/remote/portal"
	"github.com/tal-tech/go-zero/example/tracing/remote/user"
	"github.com/tal-tech/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/config.json", "the config file")

type (
	Config struct {
		zrpc.RpcServerConf
		UserRpc zrpc.RpcClientConf
	}

	PortalServer struct {
		userRpc zrpc.Client
	}
)

func NewPortalServer(client zrpc.Client) *PortalServer {
	return &PortalServer{
		userRpc: client,
	}
}

func (gs *PortalServer) Portal(ctx context.Context, req *portal.PortalRequest) (*portal.PortalResponse, error) {
	conn := gs.userRpc.Conn()
	greet := user.NewUserClient(conn)
	resp, err := greet.GetGrade(ctx, &user.UserRequest{
		Name: req.Name,
	})
	if err != nil {
		return &portal.PortalResponse{
			Response: err.Error(),
		}, nil
	} else {
		return &portal.PortalResponse{
			Response: resp.Response,
		}, nil
	}
}

func main() {
	flag.Parse()

	var c Config
	conf.MustLoad(*configFile, &c)

	client := zrpc.MustNewClient(c.UserRpc)
	server := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		portal.RegisterPortalServer(grpcServer, NewPortalServer(client))
	})
	server.Start()
}
