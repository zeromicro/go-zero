package internal

import (
	"zero/core/hash"
	"zero/core/logx"
)

type consistentBalancer struct {
	*baseBalancer
	conns     map[string]interface{}
	buckets   *hash.ConsistentHash
	bucketKey func(KV) string
}

func NewConsistentBalancer(dialFn DialFn, closeFn CloseFn, keyer func(kv KV) string) *consistentBalancer {
	// we don't support exclusive mode for consistent Balancer, to avoid complexity,
	// because there are few scenarios, use it on your own risks.
	balancer := &consistentBalancer{
		conns:     make(map[string]interface{}),
		buckets:   hash.NewConsistentHash(),
		bucketKey: keyer,
	}
	balancer.baseBalancer = newBaseBalancer(dialFn, closeFn, false)
	return balancer
}

func (b *consistentBalancer) AddConn(kv KV) error {
	// not adding kv and conn within a transaction, but it doesn't matter
	// we just rollback the kv addition if dial failed
	var conn interface{}
	prev, found := b.addKv(kv.Key, kv.Val)
	if found {
		conn = b.handlePrevious(prev)
	}

	if conn == nil {
		var err error
		conn, err = b.dialFn(kv.Val)
		if err != nil {
			b.removeKv(kv.Key)
			return err
		}
	}

	bucketKey := b.bucketKey(kv)
	b.lock.Lock()
	defer b.lock.Unlock()
	b.conns[bucketKey] = conn
	b.buckets.Add(bucketKey)
	b.notify(bucketKey)

	logx.Infof("added server, key: %s, server: %s", bucketKey, kv.Val)

	return nil
}

func (b *consistentBalancer) getConn(key string) (interface{}, bool) {
	b.lock.Lock()
	conn, ok := b.conns[key]
	b.lock.Unlock()

	return conn, ok
}

func (b *consistentBalancer) handlePrevious(prev []string) interface{} {
	if len(prev) == 0 {
		return nil
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	// if not exclusive, only need to randomly find one connection
	for key, conn := range b.conns {
		if key == prev[0] {
			return conn
		}
	}

	return nil
}

func (b *consistentBalancer) initialize() {
}

func (b *consistentBalancer) notify(key string) {
	if b.listener == nil {
		return
	}

	var keys []string
	var values []string
	for k := range b.conns {
		keys = append(keys, k)
	}
	for _, v := range b.mapping {
		values = append(values, v)
	}

	b.listener.OnUpdate(keys, values, key)
}

func (b *consistentBalancer) RemoveKey(key string) {
	kv := KV{Key: key}
	server, keep := b.removeKv(key)
	kv.Val = server
	bucketKey := b.bucketKey(kv)
	b.buckets.Remove(b.bucketKey(kv))

	// wrap the query & removal in a function to make sure the quick lock/unlock
	conn, ok := func() (interface{}, bool) {
		b.lock.Lock()
		defer b.lock.Unlock()

		conn, ok := b.conns[bucketKey]
		if ok {
			delete(b.conns, bucketKey)
		}

		return conn, ok
	}()
	if ok && !keep {
		logx.Infof("removing server, key: %s", kv.Key)
		if err := b.closeFn(server, conn); err != nil {
			logx.Error(err)
		}
	}

	// notify without new key
	b.notify("")
}

func (b *consistentBalancer) IsEmpty() bool {
	b.lock.Lock()
	empty := len(b.conns) == 0
	b.lock.Unlock()

	return empty
}

func (b *consistentBalancer) Next(keys ...string) (interface{}, bool) {
	if len(keys) != 1 {
		return nil, false
	}

	key := keys[0]
	if node, ok := b.buckets.Get(key); !ok {
		return nil, false
	} else {
		return b.getConn(node.(string))
	}
}
