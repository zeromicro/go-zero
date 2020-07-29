package internal

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/connectivity"
)

type DirectClient struct {
	conn *grpc.ClientConn
}

func NewDirectClient(server string, opts ...ClientOption) (*DirectClient, error) {
	opts = append(opts, WithDialOption(grpc.WithBalancerName(roundrobin.Name)))
	conn, err := dial(server, opts...)
	if err != nil {
		return nil, err
	}

	return &DirectClient{
		conn: conn,
	}, nil
}

func (c *DirectClient) Next() (*grpc.ClientConn, bool) {
	state := c.conn.GetState()
	if state == connectivity.Ready {
		return c.conn, true
	} else {
		return nil, false
	}
}
