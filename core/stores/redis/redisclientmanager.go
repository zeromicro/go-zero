package redis

import (
	"crypto/tls"
	"io"

	red "github.com/go-redis/redis"
	"github.com/tal-tech/go-zero/core/syncx"
)

const (
	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
)

var clientManager = syncx.NewResourceManager()

func getClient(server, pass string) (*red.Client, error) {
	return getClientWithTLS(server, pass, false)
}

func getClientWithTLS(server, pass string, tlsFlag bool) (*red.Client, error) {
	val, err := clientManager.GetResource(server, func() (io.Closer, error) {
		var tlsConfig *tls.Config = nil
		if tlsFlag {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClient(&red.Options{
			Addr:         server,
			Password:     pass,
			DB:           defaultDatabase,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
			TLSConfig:    tlsConfig,
		})
		store.WrapProcess(process)
		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}
