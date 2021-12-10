package discov

import "errors"

// EtcdConf is the config item with the given key on etcd.
type EtcdConf struct {
	Hosts []string
	Key   string
	User  string `json:",optional"`
	Pass  string `json:",optional"`
}

// HasAccount returns if account provided.
func (c EtcdConf) HasAccount() bool {
	return len(c.User) > 0 && len(c.Pass) > 0
}

// Validate validates c.
func (c EtcdConf) Validate() error {
	if len(c.Hosts) == 0 {
		return errors.New("empty etcd hosts")
	} else if len(c.Key) == 0 {
		return errors.New("empty etcd key")
	} else {
		return nil
	}
}
