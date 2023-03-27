package redis

import (
	"crypto/tls"
	"io"

	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/syncx"
)

const (
	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
)

type ClientManager interface {
	GetClient(r *Redis) (*red.Client, error)
}

type clientManager struct {
	*syncx.ResourceManager
	Database int
}

var defaultClientManager = NewClientManager(defaultDatabase)

func NewClientManager(database int) *clientManager {
	return &clientManager{
		ResourceManager: syncx.NewResourceManager(),
		Database:        database,
	}
}

func (c *clientManager) GetClient(r *Redis) (*red.Client, error) {
	val, err := c.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClient(&red.Options{
			Addr:         r.Addr,
			Password:     r.Pass,
			DB:           c.Database,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
			TLSConfig:    tlsConfig,
		})
		store.AddHook(durationHook)

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}
