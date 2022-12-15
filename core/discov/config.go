package discov

import "errors"

var (
	// errEmptyEtcdHosts indicates that etcd hosts are empty.
	errEmptyEtcdHosts = errors.New("empty etcd hosts")
	// errEmptyEtcdKey indicates that etcd key is empty.
	errEmptyEtcdKey = errors.New("empty etcd key")
)

type (
	// EtcdConf is the config item with the given key on etcd.
	EtcdConf struct {
		Hosts              []string
		Key                string
		User               string     `json:",optional"`
		Pass               string     `json:",optional"`
		CertFile           string     `json:",optional"`
		CertKeyFile        string     `json:",optional=CertFile"`
		CACertFile         string     `json:",optional=CertFile"`
		InsecureSkipVerify bool       `json:",optional"`
		Metadata           []Metadata `json:",optional"`
	}
	// Metadata represents the metadata of the service.
	Metadata struct {
		Key   string
		Value []string
	}
)

// HasAccount returns if account provided.
func (c EtcdConf) HasAccount() bool {
	return len(c.User) > 0 && len(c.Pass) > 0
}

// HasTLS returns if TLS CertFile/CertKeyFile/CACertFile are provided.
func (c EtcdConf) HasTLS() bool {
	return len(c.CertFile) > 0 && len(c.CertKeyFile) > 0 && len(c.CACertFile) > 0
}

// HasMetadata return ture if Colors exists
func (c EtcdConf) HasMetadata() bool {
	return len(c.Metadata) != 0
}

// Validate validates c.
func (c EtcdConf) Validate() error {
	if len(c.Hosts) == 0 {
		return errEmptyEtcdHosts
	} else if len(c.Key) == 0 {
		return errEmptyEtcdKey
	} else {
		return nil
	}
}
