package internal

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"sync"
	"time"
)

var (
	accounts         = make(map[string]Account)
	tlsConfigs       = make(map[string]*tls.Config)
	autoSyncIntervals = make(map[string]time.Duration)
	lock             sync.RWMutex
)

// Account holds the username/password for an etcd cluster.
type Account struct {
	User string
	Pass string
}

// AddAccount adds the username/password for the given etcd cluster.
func AddAccount(endpoints []string, user, pass string) {
	lock.Lock()
	defer lock.Unlock()

	accounts[getClusterKey(endpoints)] = Account{
		User: user,
		Pass: pass,
	}
}

// AddTLS adds the tls cert files for the given etcd cluster.
func AddTLS(endpoints []string, certFile, certKeyFile, caFile string, insecureSkipVerify bool) error {
	cert, err := tls.LoadX509KeyPair(certFile, certKeyFile)
	if err != nil {
		return err
	}

	caData, err := os.ReadFile(caFile)
	if err != nil {
		return err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	lock.Lock()
	defer lock.Unlock()
	tlsConfigs[getClusterKey(endpoints)] = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            pool,
		InsecureSkipVerify: insecureSkipVerify,
	}

	return nil
}

// GetAccount gets the username/password for the given etcd cluster.
func GetAccount(endpoints []string) (Account, bool) {
	lock.RLock()
	defer lock.RUnlock()

	account, ok := accounts[getClusterKey(endpoints)]
	return account, ok
}

// AddAutoSyncInterval adds the auto sync interval for the given etcd cluster.
func AddAutoSyncInterval(endpoints []string, interval time.Duration) {
	lock.Lock()
	defer lock.Unlock()

	autoSyncIntervals[getClusterKey(endpoints)] = interval
}

// GetAutoSyncInterval gets the auto sync interval for the given etcd cluster.
func GetAutoSyncInterval(endpoints []string) time.Duration {
	lock.RLock()
	defer lock.RUnlock()

	return autoSyncIntervals[getClusterKey(endpoints)]
}

// GetTLS gets the tls config for the given etcd cluster.
func GetTLS(endpoints []string) (*tls.Config, bool) {
	lock.RLock()
	defer lock.RUnlock()

	cfg, ok := tlsConfigs[getClusterKey(endpoints)]
	return cfg, ok
}
