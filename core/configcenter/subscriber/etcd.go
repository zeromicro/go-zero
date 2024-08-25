package subscriber

import (
	"log"

	"github.com/zeromicro/go-zero/core/discov"
)

type (
	// EtcdSubscriber is a subscriber that subscribes to etcd.
	EtcdSubscriber struct {
		*discov.Subscriber
	}

	// EtcdConfig is the configuration for etcd.
	EtcdConfig struct {
		discov.EtcdConf
	}
)

// MustNewEtcdSubscriber returns an etcd Subscriber, exits on errors.
func MustNewEtcdSubscriber(conf EtcdConfig) Subscriber {
	s, err := NewEtcdSubscriber(conf)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

// NewEtcdSubscriber returns an etcd Subscriber.
func NewEtcdSubscriber(conf EtcdConfig) (Subscriber, error) {
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

	s, err := discov.NewSubscriber(conf.EtcdConf.Hosts, conf.EtcdConf.Key, opts...)
	if err != nil {
		return nil, err
	}

	return &EtcdSubscriber{Subscriber: s}, nil
}

// AddListener adds a listener to the subscriber.
func (s *EtcdSubscriber) AddListener(listener func()) error {
	s.Subscriber.AddListener(listener)
	return nil
}

// Value returns the value of the subscriber.
func (s *EtcdSubscriber) Value() (string, error) {
	vs := s.Subscriber.Values()
	if len(vs) != 0 {
		return vs[len(vs)-1], nil
	}
	return "", nil
}
