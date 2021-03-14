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

func getClient(server, pass string, tlsFlag bool) (*red.Client, error) {
	val, err := clientManager.GetResource(server, func() (io.Closer, error) {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		if tlsFlag == false {
			tlsConfig = nil
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
