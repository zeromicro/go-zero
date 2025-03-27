package redis

import (
	"crypto/tls"
	"io"
	"runtime"

	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/syncx"
)

var (
	pubSubClientManager = syncx.NewResourceManager()
	pubSubNodePoolSize  = 10 * runtime.GOMAXPROCS(0)
)

func getPubSubClient(r *Redis) (*red.Client, error) {
	val, err := pubSubClientManager.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClient(&red.Options{
			Addr:         r.Addr,
			Username:     r.User,
			Password:     r.Pass,
			DB:           defaultDatabase,
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
			clientType: NodeType,
			key:        r.Addr,
			poolSize:   pubSubNodePoolSize,
			poolStats: func() *red.PoolStats {
				return store.PoolStats()
			},
		})

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}
