package internal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/zrpc/internal/balancer/p2c"
	"github.com/zeromicro/go-zero/zrpc/internal/clientinterceptors"
	"github.com/zeromicro/go-zero/zrpc/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	dialTimeout = time.Second * 3
	separator   = '/'
)

func init() {
	resolver.Register()
}

type (
	// Client interface wraps the Conn method.
	Client interface {
		Conn() *grpc.ClientConn
	}

	// A ClientOptions is a client options.
	ClientOptions struct {
		NonBlock    bool
		Timeout     time.Duration
		Secure      bool
		DialOptions []grpc.DialOption
	}

	// ClientOption defines the method to customize a ClientOptions.
	ClientOption func(options *ClientOptions)

	client struct {
		conn        *grpc.ClientConn
		middlewares ClientMiddlewaresConf
	}
)

// NewClient returns a Client.
func NewClient(target string, middlewares ClientMiddlewaresConf, opts ...ClientOption) (Client, error) {
	cli := client{
		middlewares: middlewares,
	}

	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, p2c.Name)
	balancerOpt := WithDialOption(grpc.WithDefaultServiceConfig(svcCfg))
	opts = append([]ClientOption{balancerOpt}, opts...)
	if err := cli.dial(target, opts...); err != nil {
		return nil, err
	}

	return &cli, nil
}

func (c *client) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *client) buildDialOptions(opts ...ClientOption) []grpc.DialOption {
	var cliOpts ClientOptions
	for _, opt := range opts {
		opt(&cliOpts)
	}

	var options []grpc.DialOption
	if !cliOpts.Secure {
		options = append([]grpc.DialOption(nil),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if !cliOpts.NonBlock {
		options = append(options, grpc.WithBlock())
	}

	options = append(options,
		grpc.WithChainUnaryInterceptor(c.buildUnaryInterceptors(cliOpts.Timeout)...),
		grpc.WithChainStreamInterceptor(c.buildStreamInterceptors()...),
	)

	return append(options, cliOpts.DialOptions...)
}

func (c *client) buildStreamInterceptors() []grpc.StreamClientInterceptor {
	var interceptors []grpc.StreamClientInterceptor

	if c.middlewares.Trace {
		interceptors = append(interceptors, clientinterceptors.StreamTracingInterceptor)
	}

	return interceptors
}

func (c *client) buildUnaryInterceptors(timeout time.Duration) []grpc.UnaryClientInterceptor {
	var interceptors []grpc.UnaryClientInterceptor

	if c.middlewares.Trace {
		interceptors = append(interceptors, clientinterceptors.UnaryTracingInterceptor)
	}
	if c.middlewares.Duration {
		interceptors = append(interceptors, clientinterceptors.DurationInterceptor)
	}
	if c.middlewares.Prometheus {
		interceptors = append(interceptors, clientinterceptors.PrometheusInterceptor)
	}
	if c.middlewares.Breaker {
		interceptors = append(interceptors, clientinterceptors.BreakerInterceptor)
	}
	if c.middlewares.Timeout {
		interceptors = append(interceptors, clientinterceptors.TimeoutInterceptor(timeout))
	}

	return interceptors
}

func (c *client) dial(server string, opts ...ClientOption) error {
	options := c.buildDialOptions(opts...)
	timeCtx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	conn, err := grpc.DialContext(timeCtx, server, options...)
	if err != nil {
		service := server
		if errors.Is(err, context.DeadlineExceeded) {
			pos := strings.LastIndexByte(server, separator)
			// len(server) - 1 is the index of last char
			if 0 < pos && pos < len(server)-1 {
				service = server[pos+1:]
			}
		}
		return fmt.Errorf("rpc dial: %s, error: %s, make sure rpc service %q is already started",
			server, err.Error(), service)
	}

	c.conn = conn
	return nil
}

// WithDialOption returns a func to customize a ClientOptions with given dial option.
func WithDialOption(opt grpc.DialOption) ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions, opt)
	}
}

// WithNonBlock sets the dialing to be nonblock.
func WithNonBlock() ClientOption {
	return func(options *ClientOptions) {
		options.NonBlock = true
	}
}

// WithStreamClientInterceptor returns a func to customize a ClientOptions with given interceptor.
func WithStreamClientInterceptor(interceptor grpc.StreamClientInterceptor) ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions,
			grpc.WithChainStreamInterceptor(interceptor))
	}
}

// WithTimeout returns a func to customize a ClientOptions with given timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(options *ClientOptions) {
		options.Timeout = timeout
	}
}

// WithTransportCredentials return a func to make the gRPC calls secured with given credentials.
func WithTransportCredentials(creds credentials.TransportCredentials) ClientOption {
	return func(options *ClientOptions) {
		options.Secure = true
		options.DialOptions = append(options.DialOptions, grpc.WithTransportCredentials(creds))
	}
}

// WithUnaryClientInterceptor returns a func to customize a ClientOptions with given interceptor.
func WithUnaryClientInterceptor(interceptor grpc.UnaryClientInterceptor) ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions,
			grpc.WithChainUnaryInterceptor(interceptor))
	}
}
