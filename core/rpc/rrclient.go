package rpc

import (
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type RRClient struct {
	conns []*grpc.ClientConn
	index int
	lock  sync.Mutex
}

func NewRRClient(endpoints []string) (*RRClient, error) {
	var conns []*grpc.ClientConn
	for _, endpoint := range endpoints {
		conn, err := dial(endpoint)
		if err != nil {
			return nil, err
		}

		conns = append(conns, conn)
	}

	rand.Seed(time.Now().UnixNano())
	return &RRClient{
		conns: conns,
		index: rand.Intn(len(conns)),
	}, nil
}

func (c *RRClient) Next() *grpc.ClientConn {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.index = (c.index + 1) % len(c.conns)
	return c.conns[c.index]
}
