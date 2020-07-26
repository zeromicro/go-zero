package discov

import (
	"sync"

	"zero/core/discov/internal"
)

type (
	subOptions struct {
		exclusive bool
	}

	SubOption func(opts *subOptions)

	Subscriber struct {
		items *container
	}
)

func NewSubscriber(endpoints []string, key string, opts ...SubOption) *Subscriber {
	var subOpts subOptions
	for _, opt := range opts {
		opt(&subOpts)
	}

	subscriber := &Subscriber{
		items: newContainer(subOpts.exclusive),
	}
	internal.GetRegistry().Monitor(endpoints, key, subscriber.items)

	return subscriber
}

func (s *Subscriber) Values() []string {
	return s.items.getValues()
}

// exclusive means that key value can only be 1:1,
// which means later added value will remove the keys associated with the same value previously.
func Exclusive() SubOption {
	return func(opts *subOptions) {
		opts.exclusive = true
	}
}

type container struct {
	exclusive bool
	values    map[string][]string
	mapping   map[string]string
	lock      sync.Mutex
}

func newContainer(exclusive bool) *container {
	return &container{
		exclusive: exclusive,
		values:    make(map[string][]string),
		mapping:   make(map[string]string),
	}
}

func (c *container) OnAdd(kv internal.KV) {
	c.addKv(kv.Key, kv.Val)
}

func (c *container) OnDelete(kv internal.KV) {
	c.removeKey(kv.Key)
}

// addKv adds the kv, returns if there are already other keys associate with the value
func (c *container) addKv(key, value string) ([]string, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

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
	} else {
		return nil, false
	}
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
	c.lock.Lock()
	defer c.lock.Unlock()

	var vs []string
	for each := range c.values {
		vs = append(vs, each)
	}
	return vs
}

// removeKey removes the kv, returns true if there are still other keys associate with the value
func (c *container) removeKey(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.doRemoveKey(key)
}

func (c *container) removeVal(val string) (empty bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for k := range c.values {
		if k == val {
			delete(c.values, k)
		}
	}
	for k, v := range c.mapping {
		if v == val {
			delete(c.mapping, k)
		}
	}

	return len(c.values) == 0
}
