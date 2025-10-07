package redis

import (
	"fmt"

	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// ClosableNode interface represents a closable redis node.
type ClosableNode interface {
	RedisNode
	Close()
}

// CreateBlockingNode creates a dedicated RedisNode for blocking operations.
//
// Blocking Redis commands (like BLPOP, BRPOP, XREADGROUP with block parameter) hold connections
// for extended periods while waiting for data. Using them with the regular Redis connection pool
// can exhaust all available connections, causing other operations to fail or timeout.
//
// CreateBlockingNode creates a separate Redis client with a minimal connection pool (size 1) that
// is dedicated to blocking operations. This ensures blocking commands don't interfere with regular
// Redis operations.
//
// Example usage:
//
//	rds := redis.MustNewRedis(redis.RedisConf{
//	    Host: "localhost:6379",
//	    Type: redis.NodeType,
//	})
//
//	// Create a dedicated node for blocking operations
//	node, err := redis.CreateBlockingNode(rds)
//	if err != nil {
//	    // handle error
//	}
//	defer node.Close() // Important: close the node when done
//
//	// Use the node for blocking operations
//	value, err := rds.Blpop(node, "mylist")
//	if err != nil {
//	    // handle error
//	}
//
// The returned ClosableNode must be closed when no longer needed to release resources.
func CreateBlockingNode(r *Redis) (ClosableNode, error) {
	timeout := readWriteTimeout + blockingQueryTimeout

	switch r.Type {
	case NodeType:
		client := red.NewClient(&red.Options{
			Addr:         r.Addr,
			Username:     r.User,
			Password:     r.Pass,
			DB:           defaultDatabase,
			MaxRetries:   maxRetries,
			PoolSize:     1,
			MinIdleConns: 1,
			ReadTimeout:  timeout,
		})
		return &clientBridge{client}, nil
	case ClusterType:
		client := red.NewClusterClient(&red.ClusterOptions{
			Addrs:        splitClusterAddrs(r.Addr),
			Username:     r.User,
			Password:     r.Pass,
			MaxRetries:   maxRetries,
			PoolSize:     1,
			MinIdleConns: 1,
			ReadTimeout:  timeout,
		})
		return &clusterBridge{client}, nil
	default:
		return nil, fmt.Errorf("unknown redis type: %s", r.Type)
	}
}

type (
	clientBridge struct {
		*red.Client
	}

	clusterBridge struct {
		*red.ClusterClient
	}
)

func (bridge *clientBridge) Close() {
	if err := bridge.Client.Close(); err != nil {
		logx.Errorf("Error occurred on close redis client: %s", err)
	}
}

func (bridge *clusterBridge) Close() {
	if err := bridge.ClusterClient.Close(); err != nil {
		logx.Errorf("Error occurred on close redis cluster: %s", err)
	}
}
