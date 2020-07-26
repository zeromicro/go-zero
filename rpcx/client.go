package rpcx

import (
	"log"
	"time"

	"zero/core/discov"
	"zero/core/rpc"
	"zero/rpcx/auth"

	"google.golang.org/grpc"
)

type RpcClient struct {
	client rpc.Client
}

func MustNewClient(c RpcClientConf, options ...rpc.ClientOption) *RpcClient {
	cli, err := NewClient(c, options...)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}

func NewClient(c RpcClientConf, options ...rpc.ClientOption) (*RpcClient, error) {
	var opts []rpc.ClientOption
	if c.HasCredential() {
		opts = append(opts, rpc.WithDialOption(grpc.WithPerRPCCredentials(&auth.Credential{
			App:   c.App,
			Token: c.Token,
		})))
	}
	if c.Timeout > 0 {
		opts = append(opts, rpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond))
	}
	opts = append(opts, options...)

	var client rpc.Client
	var err error
	if len(c.Server) > 0 {
		client, err = rpc.NewDirectClient(c.Server, opts...)
	} else if err = c.Etcd.Validate(); err == nil {
		client, err = rpc.NewRoundRobinRpcClient(c.Etcd.Hosts, c.Etcd.Key, opts...)
	}
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		client: client,
	}, nil
}

func NewClientNoAuth(c discov.EtcdConf) (*RpcClient, error) {
	client, err := rpc.NewRoundRobinRpcClient(c.Hosts, c.Key)
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		client: client,
	}, nil
}

func (rc *RpcClient) Next() (*grpc.ClientConn, bool) {
	return rc.client.Next()
}
