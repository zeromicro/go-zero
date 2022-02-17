package discov

import "github.com/zeromicro/go-zero/core/discov/internal"

// RegisterAccount registers the username/password to the given etcd cluster.
func RegisterAccount(endpoints []string, user, pass string) {
	internal.AddAccount(endpoints, user, pass)
}

// RegisterTLS registers the CertFile/CertKeyFile/CACertFile to the given etcd.
func RegisterTLS(endpoints []string, certFile, certKeyFile, caFile string,
	insecureSkipVerify bool) error {
	return internal.AddTLS(endpoints, certFile, certKeyFile, caFile, insecureSkipVerify)
}
