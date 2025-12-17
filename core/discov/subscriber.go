package discov

import (
	"sync"
	"sync/atomic"

	"github.com/zeromicro/go-zero/core/discov/internal"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
)

type (
	// SubOption defines the method to customize a Subscriber.
	SubOption func(sub *Subscriber)

	// A Subscriber is used to subscribe the given key on an etcd cluster.
	Subscriber struct {
		endpoints  []string
		exclusive  bool
		key        string
		exactMatch bool
		items      *container
	}
)

// NewSubscriber returns a Subscriber.
// endpoints is the hosts of the etcd cluster.
// key is the key to subscribe.
// opts are used to customize the Subscriber.
func NewSubscriber(endpoints []string, key string, opts ...SubOption) (*Subscriber, error) {
	sub := &Subscriber{
		endpoints: endpoints,
		key:       key,
	}
	for _, opt := range opts {
		opt(sub)
	}
	sub.items = newContainer(sub.exclusive)

	if err := internal.GetRegistry().Monitor(endpoints, key, sub.exactMatch, sub.items); err != nil {
		return nil, err
	}

	return sub, nil
}

// AddListener adds listener to s.
func (s *Subscriber) AddListener(listener func()) {
	s.items.addListener(listener)
}

// Close closes the subscriber.
func (s *Subscriber) Close() {
	internal.GetRegistry().Unmonitor(s.endpoints, s.key, s.exactMatch, s.items)
}

// Values returns all the subscription values.
func (s *Subscriber) Values() []string {
	return s.items.getValues()
}

// Exclusive means that key value can only be 1:1,
// which means later added value will remove the keys associated with the same value previously.
func Exclusive() SubOption {
	return func(sub *Subscriber) {
		sub.exclusive = true
	}
}

// WithExactMatch turn off querying using key prefixes.
func WithExactMatch() SubOption {
	return func(sub *Subscriber) {
		sub.exactMatch = true
	}
}

// WithSubEtcdAccount provides the etcd username/password.
func WithSubEtcdAccount(user, pass string) SubOption {
	return func(sub *Subscriber) {
		RegisterAccount(sub.endpoints, user, pass)
	}
}

// WithSubEtcdTLS provides the etcd CertFile/CertKeyFile/CACertFile.
func WithSubEtcdTLS(certFile, certKeyFile, caFile string, insecureSkipVerify bool) SubOption {
	return func(sub *Subscriber) {
		logx.Must(RegisterTLS(sub.endpoints, certFile, certKeyFile, caFile, insecureSkipVerify))
	}
}

type container struct {
	exclusive bool
	values    map[string][]string
	mapping   map[string]string
	snapshot  atomic.Value
	dirty     *syncx.AtomicBool
	listeners []func()
	lock      sync.Mutex
}

func newContainer(exclusive bool) *container {
	return &container{
		exclusive: exclusive,
		values:    make(map[string][]string),
		mapping:   make(map[string]string),
		dirty:     syncx.ForAtomicBool(true),
	}
}

func (c *container) OnAdd(kv internal.KV) {
	c.addKv(kv.Key, kv.Val)
	c.notifyChange()
}

func (c *container) OnDelete(kv internal.KV) {
	c.removeKey(kv.Key)
	c.notifyChange()
}

// addKv adds the kv, returns if there are already other keys associate with the value
func (c *container) addKv(key, value string) ([]string, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.dirty.Set(true)
	keys := c.values[value]
	previous := append([]string(nil), keys...)
	early := len(keys) > 0
	if c.exclusive && early {
		for _, each := range keys {
			c.doRemoveKey(each)
		}
	}
	c.values[value] = append(c.values[value], key)
	c.mapping[key] = value

	if early {
		return previous, true
	}

	return nil, false
}

func (c *container) addListener(listener func()) {
	c.lock.Lock()
	c.listeners = append(c.listeners, listener)
	c.lock.Unlock()
}

func (c *container) doRemoveKey(key string) {
	server, ok := c.mapping[key]
	if !ok {
		return
	}

	delete(c.mapping, key)
	keys := c.values[server]
	remain := keys[:0]

	for _, k := range keys {
		if k != key {
			remain = append(remain, k)
		}
	}

	if len(remain) > 0 {
		c.values[server] = remain
	} else {
		delete(c.values, server)
	}
}

func (c *container) getValues() []string {
	if !c.dirty.True() {
		return c.snapshot.Load().([]string)
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	var vals []string
	for each := range c.values {
		vals = append(vals, each)
	}
	c.snapshot.Store(vals)
	c.dirty.Set(false)

	return vals
}

func (c *container) notifyChange() {
	c.lock.Lock()
	listeners := append(([]func())(nil), c.listeners...)
	c.lock.Unlock()

	for _, listener := range listeners {
		listener()
	}
}

// removeKey removes the kv, returns true if there are still other keys associate with the value
func (c *container) removeKey(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.dirty.Set(true)
	c.doRemoveKey(key)
}
