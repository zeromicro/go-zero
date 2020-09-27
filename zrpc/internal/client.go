package internal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/zrpc/internal/balancer/p2c"
	"github.com/tal-tech/go-zero/zrpc/internal/clientinterceptors"
	"github.com/tal-tech/go-zero/zrpc/internal/resolver"
	"google.golang.org/grpc"
)

const (
	dialTimeout = time.Second * 3
	separator   = '/'
)

func init() {
	resolver.RegisterResolver()
}

type (
	ClientOptions struct {
		Timeout     time.Duration
		DialOptions []grpc.DialOption
	}

	ClientOption func(options *ClientOptions)

	client struct {
		conn *grpc.ClientConn
	}
)

func NewClient(target string, opts ...ClientOption) (*client, error) {
	opts = append(opts, WithDialOption(grpc.WithBalancerName(p2c.Name)))
	conn, err := dial(target, opts...)
	if err != nil {
		return nil, err
	}

	return &client{conn: conn}, nil
}

func (c *client) Conn() *grpc.ClientConn {
	return c.conn
}

func WithDialOption(opt grpc.DialOption) ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions, opt)
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(options *ClientOptions) {
		options.Timeout = timeout
	}
}

func buildDialOptions(opts ...ClientOption) []grpc.DialOption {
	var clientOptions ClientOptions
	for _, opt := range opts {
		opt(&clientOptions)
	}

	options := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		WithUnaryClientInterceptors(
			clientinterceptors.TracingInterceptor,
			clientinterceptors.DurationInterceptor,
			clientinterceptors.BreakerInterceptor,
			clientinterceptors.PrometheusInterceptor,
			clientinterceptors.TimeoutInterceptor(clientOptions.Timeout),
		),
	}

	return append(options, clientOptions.DialOptions...)
}

func dial(server string, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := buildDialOptions(opts...)
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
		return nil, fmt.Errorf("rpc dial: %s, error: %s, make sure rpc service %q is alread started",
			server, err.Error(), service)
	}

	return conn, nil
}
