package discov

import "github.com/z-micro/go-zero/core/discov/internal"

// RegisterAccount registers the username/password to the given etcd cluster.
func RegisterAccount(endpoints []string, user, pass string) {
	internal.AddAccount(endpoints, user, pass)
}
