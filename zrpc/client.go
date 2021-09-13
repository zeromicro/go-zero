package zrpc

import (
	"log"
	"time"

	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/zrpc/internal"
	"github.com/tal-tech/go-zero/zrpc/internal/auth"
	"google.golang.org/grpc"
)

var (
	// WithDialOption is an alias of internal.WithDialOption.
	WithDialOption = internal.WithDialOption
	// WithTimeout is an alias of internal.WithTimeout.
	WithTimeout = internal.WithTimeout
	// WithUnaryClientInterceptor is an alias of internal.WithUnaryClientInterceptor.
	WithUnaryClientInterceptor = internal.WithUnaryClientInterceptor
)

type (
	// Client is an alias of internal.Client.
	Client = internal.Client
	// ClientOption is an alias of internal.ClientOption.
	ClientOption = internal.ClientOption

	// A RpcClient is a rpc client.
	RpcClient struct {
		client Client
	}
)

// MustNewClient returns a Client, exits on any error.
func MustNewClient(c RpcClientConf, options ...ClientOption) Client {
	cli, err := NewClient(c, options...)
	if err != nil {
		log.Fatal(err)
	}

	return cli
}

// NewClient returns a Client.
func NewClient(c RpcClientConf, options ...ClientOption) (Client, error) {
	var opts []ClientOption
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

	var target string
	var err error
	if len(c.Endpoints) > 0 {
		target = internal.BuildDirectTarget(c.Endpoints)
	} else if len(c.Target) > 0 {
		target = c.Target
	} else {
		if err = c.Etcd.Validate(); err != nil {
			return nil, err
		}

		target = internal.BuildDiscovTarget(c.Etcd.Hosts, c.Etcd.Key)
	}

	client, err := internal.NewClient(target, c.Insecure, opts...)
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		client: client,
	}, nil
}

// NewClientNoAuth returns a Client without authentication.
func NewClientNoAuth(c discov.EtcdConf, opts ...ClientOption) (Client, error) {
	client, err := internal.NewClient(internal.BuildDiscovTarget(c.Hosts, c.Key), false, opts...)
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		client: client,
	}, nil
}

// NewClientWithTarget returns a Client with connecting to given target.
func NewClientWithTarget(target string, opts ...ClientOption) (Client, error) {
	return internal.NewClient(target, false, opts...)
}

// Conn returns the underlying grpc.ClientConn.
func (rc *RpcClient) Conn() *grpc.ClientConn {
	return rc.client.Conn()
}
