package subscriber

import (
	"sync"
	"sync/atomic"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	// etcdSubscriber is a subscriber that subscribes to etcd.
	etcdSubscriber struct {
		*discov.Subscriber
	}

	// EtcdConf is the configuration for etcd.
	EtcdConf = discov.EtcdConf
)

// MustNewEtcdSubscriber returns an etcd Subscriber, exits on errors.
func MustNewEtcdSubscriber(conf EtcdConf) Subscriber {
	s, err := NewEtcdSubscriber(conf)
	logx.Must(err)
	return s
}

// NewEtcdSubscriber returns an etcd Subscriber.
func NewEtcdSubscriber(conf EtcdConf) (Subscriber, error) {
	opts := buildSubOptions(conf)
	s, err := discov.NewSubscriber(conf.Hosts, conf.Key, opts...)
	if err != nil {
		return nil, err
	}

	return &etcdSubscriber{Subscriber: s}, nil
}

// buildSubOptions constructs the options for creating a new etcd subscriber.
func buildSubOptions(conf EtcdConf) []discov.SubOption {
	opts := []discov.SubOption{
		discov.WithExactMatch(),
		discov.WithContainer(newContainer()),
	}

	if len(conf.User) > 0 {
		opts = append(opts, discov.WithSubEtcdAccount(conf.User, conf.Pass))
	}
	if len(conf.CertFile) > 0 || len(conf.CertKeyFile) > 0 || len(conf.CACertFile) > 0 {
		opts = append(opts, discov.WithSubEtcdTLS(conf.CertFile, conf.CertKeyFile,
			conf.CACertFile, conf.InsecureSkipVerify))
	}

	return opts
}

// AddListener adds a listener to the subscriber.
func (s *etcdSubscriber) AddListener(listener func()) error {
	s.Subscriber.AddListener(listener)
	return nil
}

// Value returns the value of the subscriber.
func (s *etcdSubscriber) Value() (string, error) {
	vs := s.Subscriber.Values()
	if len(vs) > 0 {
		return vs[len(vs)-1], nil
	}

	return "", nil
}

type container struct {
	value     atomic.Value
	listeners []func()
	lock      sync.Mutex
}

func newContainer() *container {
	return &container{}
}

func (c *container) OnAdd(kv discov.KV) {
	c.value.Store([]string{kv.Val})
	c.notifyChange()
}

func (c *container) OnDelete(_ discov.KV) {
	c.value.Store([]string(nil))
	c.notifyChange()
}

func (c *container) AddListener(listener func()) {
	c.lock.Lock()
	c.listeners = append(c.listeners, listener)
	c.lock.Unlock()
}

func (c *container) GetValues() []string {
	if vals, ok := c.value.Load().([]string); ok {
		return vals
	}

	return []string(nil)
}

func (c *container) notifyChange() {
	c.lock.Lock()
	listeners := append(([]func())(nil), c.listeners...)
	c.lock.Unlock()

	for _, listener := range listeners {
		listener()
	}
}
