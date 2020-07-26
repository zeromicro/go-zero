package discov

import (
	"sync"

	"zero/core/discov/internal"
	"zero/core/logx"
)

const (
	_ = iota // keyBasedBalance, default
	idBasedBalance
)

type (
	Listener internal.Listener

	subClient struct {
		balancer  internal.Balancer
		lock      sync.Mutex
		cond      *sync.Cond
		listeners []internal.Listener
	}

	balanceOptions struct {
		balanceType int
	}

	BalanceOption func(*balanceOptions)

	RoundRobinSubClient struct {
		*subClient
	}

	ConsistentSubClient struct {
		*subClient
	}

	BatchConsistentSubClient struct {
		*ConsistentSubClient
	}
)

func NewRoundRobinSubClient(endpoints []string, key string, dialFn internal.DialFn, closeFn internal.CloseFn,
	opts ...SubOption) (*RoundRobinSubClient, error) {
	var subOpts subOptions
	for _, opt := range opts {
		opt(&subOpts)
	}

	cli, err := newSubClient(endpoints, key, internal.NewRoundRobinBalancer(dialFn, closeFn, subOpts.exclusive))
	if err != nil {
		return nil, err
	}

	return &RoundRobinSubClient{
		subClient: cli,
	}, nil
}

func NewConsistentSubClient(endpoints []string, key string, dialFn internal.DialFn,
	closeFn internal.CloseFn, opts ...BalanceOption) (*ConsistentSubClient, error) {
	var balanceOpts balanceOptions
	for _, opt := range opts {
		opt(&balanceOpts)
	}

	var keyer func(internal.KV) string
	switch balanceOpts.balanceType {
	case idBasedBalance:
		keyer = func(kv internal.KV) string {
			if id, ok := extractId(kv.Key); ok {
				return id
			} else {
				return kv.Key
			}
		}
	default:
		keyer = func(kv internal.KV) string {
			return kv.Val
		}
	}

	cli, err := newSubClient(endpoints, key, internal.NewConsistentBalancer(dialFn, closeFn, keyer))
	if err != nil {
		return nil, err
	}

	return &ConsistentSubClient{
		subClient: cli,
	}, nil
}

func NewBatchConsistentSubClient(endpoints []string, key string, dialFn internal.DialFn, closeFn internal.CloseFn,
	opts ...BalanceOption) (*BatchConsistentSubClient, error) {
	cli, err := NewConsistentSubClient(endpoints, key, dialFn, closeFn, opts...)
	if err != nil {
		return nil, err
	}

	return &BatchConsistentSubClient{
		ConsistentSubClient: cli,
	}, nil
}

func newSubClient(endpoints []string, key string, balancer internal.Balancer) (*subClient, error) {
	client := &subClient{
		balancer: balancer,
	}
	client.cond = sync.NewCond(&client.lock)
	if err := internal.GetRegistry().Monitor(endpoints, key, client); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *subClient) AddListener(listener internal.Listener) {
	c.lock.Lock()
	c.listeners = append(c.listeners, listener)
	c.lock.Unlock()
}

func (c *subClient) OnAdd(kv internal.KV) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if err := c.balancer.AddConn(kv); err != nil {
		logx.Error(err)
	} else {
		c.cond.Broadcast()
	}
}

func (c *subClient) OnDelete(kv internal.KV) {
	c.balancer.RemoveKey(kv.Key)
}

func (c *subClient) WaitForServers() {
	logx.Error("Waiting for alive servers")
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.balancer.IsEmpty() {
		c.cond.Wait()
	}
}

func (c *subClient) onAdd(keys []string, servers []string, newKey string) {
	// guarded by locked outside
	for _, listener := range c.listeners {
		listener.OnUpdate(keys, servers, newKey)
	}
}

func (c *RoundRobinSubClient) Next() (interface{}, bool) {
	return c.balancer.Next()
}

func (c *ConsistentSubClient) Next(key string) (interface{}, bool) {
	return c.balancer.Next(key)
}

func (bc *BatchConsistentSubClient) Next(keys []string) (map[interface{}][]string, bool) {
	if len(keys) == 0 {
		return nil, false
	}

	result := make(map[interface{}][]string)
	for _, key := range keys {
		dest, ok := bc.ConsistentSubClient.Next(key)
		if !ok {
			return nil, false
		}

		result[dest] = append(result[dest], key)
	}

	return result, true
}

func BalanceWithId() BalanceOption {
	return func(opts *balanceOptions) {
		opts.balanceType = idBasedBalance
	}
}
