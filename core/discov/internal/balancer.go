package internal

import "sync"

type (
	DialFn  func(server string) (interface{}, error)
	CloseFn func(server string, conn interface{}) error

	Balancer interface {
		AddConn(kv KV) error
		IsEmpty() bool
		Next(key ...string) (interface{}, bool)
		RemoveKey(key string)
		initialize()
		setListener(listener Listener)
	}

	serverConn struct {
		key  string
		conn interface{}
	}

	baseBalancer struct {
		exclusive bool
		servers   map[string][]string
		mapping   map[string]string
		lock      sync.Mutex
		dialFn    DialFn
		closeFn   CloseFn
		listener  Listener
	}
)

func newBaseBalancer(dialFn DialFn, closeFn CloseFn, exclusive bool) *baseBalancer {
	return &baseBalancer{
		exclusive: exclusive,
		servers:   make(map[string][]string),
		mapping:   make(map[string]string),
		dialFn:    dialFn,
		closeFn:   closeFn,
	}
}

// addKv adds the kv, returns if there are already other keys associate with the server
func (b *baseBalancer) addKv(key, value string) ([]string, bool) {
	b.lock.Lock()
	defer b.lock.Unlock()

	keys := b.servers[value]
	previous := append([]string(nil), keys...)
	early := len(keys) > 0
	if b.exclusive && early {
		for _, each := range keys {
			b.doRemoveKv(each)
		}
	}
	b.servers[value] = append(b.servers[value], key)
	b.mapping[key] = value

	if early {
		return previous, true
	} else {
		return nil, false
	}
}

func (b *baseBalancer) doRemoveKv(key string) (server string, keepConn bool) {
	server, ok := b.mapping[key]
	if !ok {
		return "", true
	}

	delete(b.mapping, key)
	keys := b.servers[server]
	remain := keys[:0]

	for _, k := range keys {
		if k != key {
			remain = append(remain, k)
		}
	}

	if len(remain) > 0 {
		b.servers[server] = remain
		return server, true
	} else {
		delete(b.servers, server)
		return server, false
	}
}

func (b *baseBalancer) removeKv(key string) (server string, keepConn bool) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.doRemoveKv(key)
}

func (b *baseBalancer) setListener(listener Listener) {
	b.lock.Lock()
	b.listener = listener
	b.lock.Unlock()
}
