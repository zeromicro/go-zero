package discov

import "errors"

type EtcdConf struct {
	Hosts []string
	Key   string
}

func (c EtcdConf) Validate() error {
	if len(c.Hosts) == 0 {
		return errors.New("empty etcd hosts")
	} else if len(c.Key) == 0 {
		return errors.New("empty etcd key")
	} else {
		return nil
	}
}

type K8sConf struct {
	Name      string
	Namespace string
	Port      uint32
}

func (c K8sConf) Validate() error {
	if len(c.Name) == 0 {
		return errors.New("service name empty")
	}
	if len(c.Namespace) == 0 {
		return errors.New("service namespace empty")
	}
	if c.Port == 0 {
		return errors.New("service port shoud not be zero")
	}
	return nil
}
