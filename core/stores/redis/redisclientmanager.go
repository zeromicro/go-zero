package redis

import (
	"crypto/tls"
	"io"
	"runtime"

	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/syncx"
)

const (
	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
)

var (
	clientManager = syncx.NewResourceManager()
	// nodePoolSize is default pool size for node type of redis.
	nodePoolSize = 10 * runtime.GOMAXPROCS(0)
)

func getClient(r *Redis) (*red.Client, error) {
	val, err := clientManager.GetResource(r.Addr, func() (io.Closer, error) {
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
			DB:           r.db,
			MaxRetries:   r.maxRetries,
			MinIdleConns: r.minIdleConns,
			PoolSize:     r.poolSize,
			TLSConfig:    tlsConfig,
		})

		hooks := append([]red.Hook{defaultDurationHook, breakerHook{
			brk: r.brk,
		}}, r.hooks...)
		for _, hook := range hooks {
			store.AddHook(hook)
		}

		poolSize := nodePoolSize
		if r.poolSize > 0 {
			poolSize = r.poolSize
		}

		connCollector.registerClient(&statGetter{
			clientType: NodeType,
			key:        r.Addr,
			poolSize:   poolSize,
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
