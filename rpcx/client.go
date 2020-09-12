package rpcx

import (
	"log"
	"time"

	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/rpcx/internal"
	"github.com/tal-tech/go-zero/rpcx/internal/auth"
	"google.golang.org/grpc"
)

var (
	WithDialOption = internal.WithDialOption
	WithTimeout    = internal.WithTimeout
)

type (
	Client interface {
		Conn() *grpc.ClientConn
	}

	RpcClient struct {
		client Client
	}
)

func MustNewClient(c RpcClientConf, options ...internal.ClientOption) Client {
	cli, err := NewClient(c, options...)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}

func NewClient(c RpcClientConf, options ...internal.ClientOption) (Client, error) {
	var opts []internal.ClientOption
	if c.HasCredential() {
		opts = append(opts, WithDialOption(grpc.WithPerRPCCredentials(&auth.Credential{
			App:   c.App,
			Token: c.Token,
		})))
	}
	if c.Timeout > 0 {
		opts = append(opts, WithTimeout(time.Duration(c.Timeout)*time.Millisecond))
	}
	opts = append(opts, options...)

	var client Client
	var err error
	if len(c.Endpoints) > 0 {
		client, err = internal.NewClient(internal.BuildDirectTarget(c.Endpoints), opts...)
	} else if err = c.Etcd.Validate(); err == nil {
		client, err = internal.NewClient(internal.BuildDiscovTarget(c.Etcd.Hosts, c.Etcd.Key), opts...)
	}
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		client: client,
	}, nil
}

func NewClientNoAuth(c discov.EtcdConf) (Client, error) {
	client, err := internal.NewClient(internal.BuildDiscovTarget(c.Hosts, c.Key))
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		client: client,
	}, nil
}

func NewClientWithTarget(target string, opts ...internal.ClientOption) (Client, error) {
	return internal.NewClient(target, opts...)
}

func (rc *RpcClient) Conn() *grpc.ClientConn {
	return rc.client.Conn()
}
