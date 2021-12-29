package discov

import "github.com/tal-tech/go-zero/core/discov/internal"

// RegisterAccount registers the username/password to the given etcd.
func RegisterAccount(endpoints []string, user, pass string) {
	internal.AddAccount(endpoints, user, pass)
}

// RegisterTLS registers the CertFile/CertKeyFile/TrustedCAFile to the given etcd.
func RegisterTLS(endpoints []string, certFile, certKeyFile, caFile string) error {
	return internal.AddTLS(endpoints, certFile, certKeyFile, caFile)
}
