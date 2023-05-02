package zrpc

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc/internal"
	"github.com/zeromicro/go-zero/zrpc/internal/auth"
	"github.com/zeromicro/go-zero/zrpc/internal/clientinterceptors"
)

const defaultClientKeepaliveTime = 20 * time.Second

var (
	// WithDialOption is an alias of internal.WithDialOption.
	WithDialOption = internal.WithDialOption
	// WithNonBlock sets the dialing to be nonblock.
	WithNonBlock = internal.WithNonBlock
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
	logx.Must(err)
	return cli
}

// NewClientIfEnable returns a Client if it is enabled.
func NewClientIfEnable(c RpcClientConf, options ...ClientOption) Client {
	if !c.Enabled {
		return nil
	}

	cli, err := NewClient(c, options...)
	if err != nil {
		logx.Must(err)
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
	if c.KeepaliveTime > 0 {
		opts = append(opts, WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: c.KeepaliveTime,
		})))
	}

	opts = append(opts, options...)

	target, err := c.BuildTarget()
	if err != nil {
		return nil, err
	}

	client, err := internal.NewClient(target, c.Middlewares, opts...)
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		client: client,
	}, nil
}

// NewClientWithTarget returns a Client with connecting to given target.
func NewClientWithTarget(target string, opts ...ClientOption) (Client, error) {
	middlewares := ClientMiddlewaresConf{
		Trace:      true,
		Duration:   true,
		Prometheus: true,
		Breaker:    true,
		Timeout:    true,
	}

	opts = append([]ClientOption{
		WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: defaultClientKeepaliveTime,
		})),
	}, opts...)

	return internal.NewClient(target, middlewares, opts...)
}

// Conn returns the underlying grpc.ClientConn.
func (rc *RpcClient) Conn() *grpc.ClientConn {
	return rc.client.Conn()
}

// DontLogClientContentForMethod disable logging content for given method.
func DontLogClientContentForMethod(method string) {
	clientinterceptors.DontLogContentForMethod(method)
}

// SetClientSlowThreshold sets the slow threshold on client side.
func SetClientSlowThreshold(threshold time.Duration) {
	clientinterceptors.SetSlowThreshold(threshold)
}
