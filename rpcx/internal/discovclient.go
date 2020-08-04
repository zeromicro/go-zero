package internal

import (
	"zero/core/discov"
	"zero/rpcx/internal/balancer/roundrobin"
	"zero/rpcx/internal/resolver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type DiscovClient struct {
	conn *grpc.ClientConn
}

func NewDiscovClient(etcd discov.EtcdConf, opts ...ClientOption) (*DiscovClient, error) {
	resolver.RegisterResolver(etcd)
	opts = append(opts, WithDialOption(grpc.WithBalancerName(roundrobin.Name)))
	conn, err := dial("discov:///", opts...)
	if err != nil {
		return nil, err
	}

	return &DiscovClient{
		conn: conn,
	}, nil
}

func (c *DiscovClient) Next() (*grpc.ClientConn, bool) {
	state := c.conn.GetState()
	if state == connectivity.Ready {
		return c.conn, true
	} else {
		return nil, false
	}
}
