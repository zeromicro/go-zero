package redis

import (
	"crypto/tls"
	"io"

	red "github.com/go-redis/redis"
	"github.com/tal-tech/go-zero/core/syncx"
)

var clusterManager = syncx.NewResourceManager()

func getCluster(server, pass string) (*red.ClusterClient, error) {
	return getClusterWithTLS(server, pass, false)
}

func getClusterWithTLS(server, pass string, tlsEnabled bool) (*red.ClusterClient, error) {
	val, err := clusterManager.GetResource(server, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if tlsEnabled {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClusterClient(&red.ClusterOptions{
			Addrs:        []string{server},
			Password:     pass,
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

	return val.(*red.ClusterClient), nil
}
