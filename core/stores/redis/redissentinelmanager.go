package redis

import (
	"crypto/tls"
	"io"
	"runtime"

	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/syncx"
)

var (
	sentinelManager  = syncx.NewResourceManager()
	sentinelPoolSize = 5 * runtime.GOMAXPROCS(0)
)

func getSentinel(r *Redis) (*red.Client, error) {
	val, err := sentinelManager.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewFailoverClient(&red.FailoverOptions{
			SentinelAddrs: splitClusterAddrs(r.Addr),
			MasterName:    r.MasterName,
			Username:      r.User,
			Password:      r.Pass,
			MaxRetries:    maxRetries,
			MinIdleConns:  idleConns,
			TLSConfig:     tlsConfig,
		})

		hooks := append([]red.Hook{defaultDurationHook, breakerHook{
			brk: r.brk,
		}}, r.hooks...)
		for _, hook := range hooks {
			store.AddHook(hook)
		}
		connCollector.registerClient(&statGetter{
			clientType: SentinelType,
			key:        r.Addr,
			poolSize:   sentinelPoolSize,
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
