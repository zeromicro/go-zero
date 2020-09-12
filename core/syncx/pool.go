package syncx

import (
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/timex"
)

type (
	PoolOption func(*Pool)

	node struct {
		item     interface{}
		next     *node
		lastUsed time.Duration
	}

	Pool struct {
		limit   int
		created int
		maxAge  time.Duration
		lock    sync.Locker
		cond    *sync.Cond
		head    *node
		create  func() interface{}
		destroy func(interface{})
	}
)

func NewPool(n int, create func() interface{}, destroy func(interface{}), opts ...PoolOption) *Pool {
	if n <= 0 {
		panic("pool size can't be negative or zero")
	}

	lock := new(sync.Mutex)
	pool := &Pool{
		limit:   n,
		lock:    lock,
		cond:    sync.NewCond(lock),
		create:  create,
		destroy: destroy,
	}

	for _, opt := range opts {
		opt(pool)
	}

	return pool
}

func (p *Pool) Get() interface{} {
	p.lock.Lock()
	defer p.lock.Unlock()

	for {
		if p.head != nil {
			head := p.head
			p.head = head.next
			if p.maxAge > 0 && head.lastUsed+p.maxAge < timex.Now() {
				p.created--
				p.destroy(head.item)
				continue
			} else {
				return head.item
			}
		}

		if p.created < p.limit {
			p.created++
			return p.create()
		}

		p.cond.Wait()
	}
}

func (p *Pool) Put(x interface{}) {
	if x == nil {
		return
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.head = &node{
		item:     x,
		next:     p.head,
		lastUsed: timex.Now(),
	}
	p.cond.Signal()
}

func WithMaxAge(duration time.Duration) PoolOption {
	return func(pool *Pool) {
		pool.maxAge = duration
	}
}
