package zrpc

import (
	"log"
	"time"

	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/zrpc/internal"
	"github.com/tal-tech/go-zero/zrpc/internal/auth"
	"github.com/tal-tech/go-zero/zrpc/internal/clientinterceptors"
	"google.golang.org/grpc"
)

var (
	// WithDialOption is an alias of internal.WithDialOption.
	WithDialOption = internal.WithDialOption
	// WithNonBlock sets the dialing to be nonblock.
	WithNonBlock = internal.WithNonBlock
	// WithTimeout is an alias of internal.WithTimeout.
	WithTimeout = internal.WithTimeout
	// WithRetry is an alias of internal.WithRetry.
	WithRetry = internal.WithRetry
	// WithTransportCredentials return a func to make the gRPC calls secured with given credentials.
	WithTransportCredentials = internal.WithTransportCredentials
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
	if c.NonBlock {
		opts = append(opts, WithNonBlock())
	}
	if c.Timeout > 0 {
		opts = append(opts, WithTimeout(time.Duration(c.Timeout)*time.Millisecond))
	}
	if c.Retry {
		opts = append(opts, WithRetry())
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

		if c.Etcd.HasAccount() {
			discov.RegisterAccount(c.Etcd.Hosts, c.Etcd.User, c.Etcd.Pass)
		}

		target = internal.BuildDiscovTarget(c.Etcd.Hosts, c.Etcd.Key)
	}

	client, err := internal.NewClient(target, opts...)
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		client: client,
	}, nil
}

// NewClientWithTarget returns a Client with connecting to given target.
func NewClientWithTarget(target string, opts ...ClientOption) (Client, error) {
	return internal.NewClient(target, opts...)
}

// Conn returns the underlying grpc.ClientConn.
func (rc *RpcClient) Conn() *grpc.ClientConn {
	return rc.client.Conn()
}

// SetClientSlowThreshold sets the slow threshold on client side.
func SetClientSlowThreshold(threshold time.Duration) {
	clientinterceptors.SetSlowThreshold(threshold)
}
