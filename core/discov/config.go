package discov

import "errors"

// EtcdConf is the config item with the given key on etcd.
type EtcdConf struct {
	Hosts []string
	Key   string
	User  string `json:",optional"`
	Pass  string `json:",optional=User"`
}

// Validate validates c.
func (c EtcdConf) Validate() error {
	if len(c.Hosts) == 0 {
		return errors.New("empty etcd hosts")
	} else if len(c.Key) == 0 {
		return errors.New("empty etcd key")
	} else if len(c.User) > 0 || len(c.Pass) > 0 {
		if len(c.User) == 0 || len(c.Pass) == 0 {
			return errors.New("etcd user and pass must pairs")
		}
		return nil
	} else {
		return nil
	}
}

func (c EtcdConf) EnableAuth() bool {
	return len(c.User) > 0 && len(c.Pass) > 0
}
