package discov

import (
	"time"

	"github.com/zeromicro/go-zero/core/discov/internal"
)

// RegisterAccount registers the username/password to the given etcd cluster.
func RegisterAccount(endpoints []string, user, pass string) {
	internal.AddAccount(endpoints, user, pass)
}

// RegisterTLS registers the CertFile/CertKeyFile/CACertFile to the given etcd.
func RegisterTLS(endpoints []string, certFile, certKeyFile, caFile string,
	insecureSkipVerify bool) error {
	return internal.AddTLS(endpoints, certFile, certKeyFile, caFile, insecureSkipVerify)
}

// RegisterAutoSyncInterval registers the auto sync interval to the given etcd
// cluster, a zero interval disables the auto sync of the etcd client.
func RegisterAutoSyncInterval(endpoints []string, interval time.Duration) {
	internal.AddAutoSyncInterval(endpoints, interval)
}
