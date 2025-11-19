package zrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc/internal"
	"github.com/zeromicro/go-zero/zrpc/internal/auth"
	"github.com/zeromicro/go-zero/zrpc/internal/balancer/consistenthash"
	"github.com/zeromicro/go-zero/zrpc/internal/balancer/p2c"
	"github.com/zeromicro/go-zero/zrpc/internal/clientinterceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	// WithDialOption is an alias of internal.WithDialOption.
	WithDialOption = internal.WithDialOption
	// WithNonBlock sets the dialing to be nonblock.
	WithNonBlock = internal.WithNonBlock
	// WithBlock sets the dialing to be blocking.
	// Deprecated: blocking dials are not recommended by gRPC.
	WithBlock = internal.WithBlock
	// WithStreamClientInterceptor is an alias of internal.WithStreamClientInterceptor.
	WithStreamClientInterceptor = internal.WithStreamClientInterceptor
	// WithTimeout is an alias of internal.WithTimeout.
	WithTimeout = internal.WithTimeout
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
	logx.Must(err)
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
	} else {
		opts = append(opts, WithBlock())
	}
	if c.Timeout > 0 {
		opts = append(opts, WithTimeout(time.Duration(c.Timeout)*time.Millisecond))
	}
	if c.KeepaliveTime > 0 {
		opts = append(opts, WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: c.KeepaliveTime,
		})))
	}

	svcCfg := makeLBServiceConfig(c.BalancerName)
	opts = append(opts, WithDialOption(grpc.WithDefaultServiceConfig(svcCfg)))

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
	var config RpcClientConf
	if err := conf.FillDefault(&config); err != nil {
		return nil, err
	}

	config.Target = target

	return NewClient(config, opts...)
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

// SetHashKey sets the hash key into context.
func SetHashKey(ctx context.Context, key string) context.Context {
	return consistenthash.SetHashKey(ctx, key)
}

// WithCallTimeout return a call option with given timeout to make a method call.
func WithCallTimeout(timeout time.Duration) grpc.CallOption {
	return clientinterceptors.WithCallTimeout(timeout)
}

func makeLBServiceConfig(balancerName string) string {
	if len(balancerName) == 0 {
		balancerName = p2c.Name
	}

	return fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, balancerName)
}
