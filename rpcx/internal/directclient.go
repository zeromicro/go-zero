package internal

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
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

func (c *DirectClient) Conn() *grpc.ClientConn {
	return c.conn
}
