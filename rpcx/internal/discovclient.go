package internal

import (
	"fmt"
	"strings"

	"zero/rpcx/internal/balancer/p2c"
	"zero/rpcx/internal/resolver"

	"google.golang.org/grpc"
)

func init() {
	resolver.RegisterResolver()
}

type DiscovClient struct {
	conn *grpc.ClientConn
}

func NewDiscovClient(endpoints []string, key string, opts ...ClientOption) (*DiscovClient, error) {
	opts = append(opts, WithDialOption(grpc.WithBalancerName(p2c.Name)))
	target := fmt.Sprintf("%s://%s/%s", resolver.DiscovScheme,
		strings.Join(endpoints, resolver.EndpointSep), key)
	conn, err := dial(target, opts...)
	if err != nil {
		return nil, err
	}

	return &DiscovClient{conn: conn}, nil
}

func (c *DiscovClient) Conn() *grpc.ClientConn {
	return c.conn
}
