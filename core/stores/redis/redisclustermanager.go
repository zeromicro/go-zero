package redis

import (
	"crypto/tls"
	"io"
	"strings"

	red "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/syncx"
)

const addrSep = ","

var clusterManager = syncx.NewResourceManager()

func getCluster(r *Redis) (*red.ClusterClient, error) {
	val, err := clusterManager.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClusterClient(&red.ClusterOptions{
			Addrs:        splitClusterAddrs(r.Addr),
			Password:     r.Pass,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
			TLSConfig:    tlsConfig,
		})
		store.AddHook(durationHook)
		for _, hook := range r.hooks {
			store.AddHook(hook)
		}

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.ClusterClient), nil
}

func splitClusterAddrs(addr string) []string {
	addrs := strings.Split(addr, addrSep)
	unique := make(map[string]struct{})
	for _, each := range addrs {
		unique[strings.TrimSpace(each)] = struct{}{}
	}

	addrs = addrs[:0]
	for k := range unique {
		addrs = append(addrs, k)
	}

	return addrs
}
