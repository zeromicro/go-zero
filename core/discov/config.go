package discov

import "errors"

var (
	// errEmptyEtcdHosts indicates that etcd hosts are empty.
	errEmptyEtcdHosts = errors.New("empty etcd hosts")
	// errEmptyEtcdKey indicates that etcd key is empty.
	errEmptyEtcdKey = errors.New("empty etcd key")
)

// EtcdConf is the config item with the given key on etcd.
type EtcdConf struct {
	Hosts              []string `json:"Hosts" yaml:"Hosts"`
	Key                string   `json:"Key" yaml:"Key"`
	User               string   `json:"User,optional" yaml:"User"`
	Pass               string   `json:"Pass,optional" yaml:"Pass"`
	CertFile           string   `json:"CertFile,optional" yaml:"CertFile"`
	CertKeyFile        string   `json:"CertKeyFile,optional=CertFile" yaml:"CertKeyFile"`
	CACertFile         string   `json:"CACertFile,optional=CertFile" yaml:"CACertFile"`
	InsecureSkipVerify bool     `json:"InsecureSkipVerify,optional" yaml:"InsecureSkipVerify"`
}

// HasAccount returns if account provided.
func (c EtcdConf) HasAccount() bool {
	return len(c.User) > 0 && len(c.Pass) > 0
}

// HasTLS returns if TLS CertFile/CertKeyFile/CACertFile are provided.
func (c EtcdConf) HasTLS() bool {
	return len(c.CertFile) > 0 && len(c.CertKeyFile) > 0 && len(c.CACertFile) > 0
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
