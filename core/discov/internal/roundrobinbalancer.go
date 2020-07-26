package internal

import (
	"math/rand"
	"time"

	"zero/core/logx"
)

type roundRobinBalancer struct {
	*baseBalancer
	conns []serverConn
	index int
}

func NewRoundRobinBalancer(dialFn DialFn, closeFn CloseFn, exclusive bool) *roundRobinBalancer {
	balancer := new(roundRobinBalancer)
	balancer.baseBalancer = newBaseBalancer(dialFn, closeFn, exclusive)
	return balancer
}

func (b *roundRobinBalancer) AddConn(kv KV) error {
	var conn interface{}
	prev, found := b.addKv(kv.Key, kv.Val)
	if found {
		conn = b.handlePrevious(prev, kv.Val)
	}

	if conn == nil {
		var err error
		conn, err = b.dialFn(kv.Val)
		if err != nil {
			b.removeKv(kv.Key)
			return err
		}
	}

	b.lock.Lock()
	defer b.lock.Unlock()
	b.conns = append(b.conns, serverConn{
		key:  kv.Key,
		conn: conn,
	})
	b.notify(kv.Key)

	return nil
}

func (b *roundRobinBalancer) handlePrevious(prev []string, server string) interface{} {
	if len(prev) == 0 {
		return nil
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	if b.exclusive {
		for _, item := range prev {
			conns := b.conns[:0]
			for _, each := range b.conns {
				if each.key == item {
					if err := b.closeFn(server, each.conn); err != nil {
						logx.Error(err)
					}
				} else {
					conns = append(conns, each)
				}
			}
			b.conns = conns
		}
	} else {
		for _, each := range b.conns {
			if each.key == prev[0] {
				return each.conn
			}
		}
	}

	return nil
}

func (b *roundRobinBalancer) initialize() {
	rand.Seed(time.Now().UnixNano())
	if len(b.conns) > 0 {
		b.index = rand.Intn(len(b.conns))
	}
}

func (b *roundRobinBalancer) IsEmpty() bool {
	b.lock.Lock()
	empty := len(b.conns) == 0
	b.lock.Unlock()

	return empty
}

func (b *roundRobinBalancer) Next(...string) (interface{}, bool) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if len(b.conns) == 0 {
		return nil, false
	}

	b.index = (b.index + 1) % len(b.conns)
	return b.conns[b.index].conn, true
}

func (b *roundRobinBalancer) notify(key string) {
	if b.listener == nil {
		return
	}

	// b.servers has the format of map[conn][]key
	var keys []string
	var values []string
	for k, v := range b.servers {
		values = append(values, k)
		keys = append(keys, v...)
	}

	b.listener.OnUpdate(keys, values, key)
}

func (b *roundRobinBalancer) RemoveKey(key string) {
	server, keep := b.removeKv(key)

	b.lock.Lock()
	defer b.lock.Unlock()

	conns := b.conns[:0]
	for _, conn := range b.conns {
		if conn.key == key {
			// there are other keys assocated with the conn, don't close the conn.
			if keep {
				continue
			}
			if err := b.closeFn(server, conn.conn); err != nil {
				logx.Error(err)
			}
		} else {
			conns = append(conns, conn)
		}
	}
	b.conns = conns
	// notify without new key
	b.notify("")
}
