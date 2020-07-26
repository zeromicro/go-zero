package rpc

import (
	"time"

	"zero/core/discov"
	"zero/core/logx"
	"zero/core/threading"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	coolOffTime = time.Second * 5
	retryTimes  = 3
)

type (
	RoundRobinSubClient struct {
		*discov.RoundRobinSubClient
	}

	ConsistentSubClient struct {
		*discov.ConsistentSubClient
	}
)

func NewRoundRobinRpcClient(endpoints []string, key string, opts ...ClientOption) (*RoundRobinSubClient, error) {
	subClient, err := discov.NewRoundRobinSubClient(endpoints, key, func(server string) (interface{}, error) {
		return dial(server, opts...)
	}, func(server string, conn interface{}) error {
		return closeConn(conn.(*grpc.ClientConn))
	}, discov.Exclusive())
	if err != nil {
		return nil, err
	} else {
		return &RoundRobinSubClient{subClient}, nil
	}
}

func NewConsistentRpcClient(endpoints []string, key string, opts ...ClientOption) (*ConsistentSubClient, error) {
	subClient, err := discov.NewConsistentSubClient(endpoints, key, func(server string) (interface{}, error) {
		return dial(server, opts...)
	}, func(server string, conn interface{}) error {
		return closeConn(conn.(*grpc.ClientConn))
	})
	if err != nil {
		return nil, err
	} else {
		return &ConsistentSubClient{subClient}, nil
	}
}

func (cli *RoundRobinSubClient) Next() (*grpc.ClientConn, bool) {
	return next(func() (interface{}, bool) {
		return cli.RoundRobinSubClient.Next()
	})
}

func (cli *ConsistentSubClient) Next(key string) (*grpc.ClientConn, bool) {
	return next(func() (interface{}, bool) {
		return cli.ConsistentSubClient.Next(key)
	})
}

func closeConn(conn *grpc.ClientConn) error {
	// why to close the conn asynchronously is because maybe another goroutine
	// is using the same conn, we can wait the coolOffTime to let the other
	// goroutine to finish using the conn.
	// after the conn unregistered, the balancer will not assign the conn,
	// but maybe the already assigned tasks are still using it.
	threading.GoSafe(func() {
		time.Sleep(coolOffTime)
		if err := conn.Close(); err != nil {
			logx.Error(err)
		}
	})

	return nil
}

func next(nextFn func() (interface{}, bool)) (*grpc.ClientConn, bool) {
	for i := 0; i < retryTimes; i++ {
		v, ok := nextFn()
		if !ok {
			break
		}

		conn, yes := v.(*grpc.ClientConn)
		if !yes {
			break
		}

		switch conn.GetState() {
		case connectivity.Ready:
			return conn, true
		}
	}

	return nil, false
}
