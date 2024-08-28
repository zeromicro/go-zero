package subscriber

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	// etcdSubscriber is a subscriber that subscribes to etcd.
	etcdSubscriber struct {
		*discov.Subscriber
	}

	// EtcdConf is the configuration for etcd.
	EtcdConf discov.EtcdConf
)

// MustNewEtcdSubscriber returns an etcd Subscriber, exits on errors.
func MustNewEtcdSubscriber(conf EtcdConf) Subscriber {
	s, err := NewEtcdSubscriber(conf)
	logx.Must(err)

	return s
}

// NewEtcdSubscriber returns an etcd Subscriber.
func NewEtcdSubscriber(conf EtcdConf) (Subscriber, error) {
	var opts = []discov.SubOption{
		discov.WithDisablePrefix(),
	}
	if len(conf.User) != 0 {
		opts = append(opts, discov.WithSubEtcdAccount(conf.User, conf.Pass))
	}
	if len(conf.CertFile) != 0 || len(conf.CertKeyFile) != 0 || len(conf.CACertFile) != 0 {
		opts = append(opts,
			discov.WithSubEtcdTLS(conf.CertFile, conf.CertKeyFile, conf.CACertFile, conf.InsecureSkipVerify))
	}

	s, err := discov.NewSubscriber(conf.Hosts, conf.Key, opts...)
	if err != nil {
		return nil, err
	}

	return &etcdSubscriber{Subscriber: s}, nil
}

// AddListener adds a listener to the subscriber.
func (s *etcdSubscriber) AddListener(listener func()) error {
	s.Subscriber.AddListener(listener)
	return nil
}

// Value returns the value of the subscriber.
func (s *etcdSubscriber) Value() (string, error) {
	vs := s.Subscriber.Values()
	if len(vs) != 0 {
		return vs[len(vs)-1], nil
	}
	return "", nil
}
