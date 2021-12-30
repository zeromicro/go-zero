package internal

import "sync"

var (
	accounts = make(map[string]Account)
	lock     sync.RWMutex
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

// GetAccount gets the username/password for the given etcd cluster.
func GetAccount(endpoints []string) (Account, bool) {
	lock.RLock()
	defer lock.RUnlock()

	account, ok := accounts[getClusterKey(endpoints)]
	return account, ok
}
