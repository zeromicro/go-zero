package redis

import (
	"crypto/tls"
	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/syncx"
	"io"
	"runtime"
)

var (
	pubSubClusterManager = syncx.NewResourceManager()
	// clusterPoolSize is default pool size for cluster type of redis.
	pubSubClusterPoolSize = 5 * runtime.GOMAXPROCS(0)
)

func getPubSubCluster(r *Redis) (*red.ClusterClient, error) {
	val, err := pubSubClusterManager.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClusterClient(&red.ClusterOptions{
			Addrs:        splitClusterAddrs(r.Addr),
			Username:     r.User,
			Password:     r.Pass,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
			TLSConfig:    tlsConfig,
		})

		hooks := append([]red.Hook{defaultDurationHook, breakerHook{
			brk: r.brk,
		}}, r.hooks...)
		for _, hook := range hooks {
			store.AddHook(hook)
		}

		connCollector.registerClient(&statGetter{
			clientType: ClusterType,
			key:        r.Addr,
			poolSize:   pubSubClusterPoolSize,
			poolStats: func() *red.PoolStats {
				return store.PoolStats()
			},
		})

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.ClusterClient), nil
}
