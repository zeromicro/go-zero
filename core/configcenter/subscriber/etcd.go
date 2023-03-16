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

// Values returns the values of the subscriber.
func (s *EtcdSubscriber) Values() ([]string, error) {
	return s.Subscriber.Values(), nil
}
