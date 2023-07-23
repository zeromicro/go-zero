package syncx

import (
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/timex"
)

type (
	// PoolOption defines the method to customize a Pool.
	PoolOption func(*Pool)

	node struct {
		item     any
		next     *node
		lastUsed time.Duration
	}

	// A Pool is used to pool resources.
	// The difference between sync.Pool is that:
	//  1. the limit of the resources
	//  2. max age of the resources can be set
	//  3. the method to destroy resources can be customized
	Pool struct {
		limit   int
		created int
		maxAge  time.Duration
		lock    sync.Locker
		cond    *sync.Cond
		head    *node
		create  func() any
		destroy func(any)
	}
)

// NewPool returns a Pool.
func NewPool(n int, create func() any, destroy func(any), opts ...PoolOption) *Pool {
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

// Get gets a resource.
func (p *Pool) Get() any {
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

// Put puts a resource back.
func (p *Pool) Put(x any) {
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

// WithMaxAge returns a function to customize a Pool with given max age.
func WithMaxAge(duration time.Duration) PoolOption {
	return func(pool *Pool) {
		pool.maxAge = duration
	}
}
