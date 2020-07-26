package syncx

import "sync"

type (
	// SharedCalls lets the concurrent calls with the same key to share the call result.
	// For example, A called F, before it's done, B called F. Then B would not execute F,
	// and shared the result returned by F which called by A.
	// The calls with the same key are dependent, concurrent calls share the returned values.
	// A ------->calls F with key<------------------->returns val
	// B --------------------->calls F with key------>returns val
	SharedCalls interface {
		Do(key string, fn func() (interface{}, error)) (interface{}, error)
		DoEx(key string, fn func() (interface{}, error)) (interface{}, bool, error)
	}

	call struct {
		wg  sync.WaitGroup
		val interface{}
		err error
	}

	sharedGroup struct {
		calls map[string]*call
		lock  sync.Mutex
	}
)

func NewSharedCalls() SharedCalls {
	return &sharedGroup{
		calls: make(map[string]*call),
	}
}

func (g *sharedGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.lock.Lock()
	if c, ok := g.calls[key]; ok {
		g.lock.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := g.makeCall(key, fn)
	return c.val, c.err
}

func (g *sharedGroup) DoEx(key string, fn func() (interface{}, error)) (val interface{}, fresh bool, err error) {
	g.lock.Lock()
	if c, ok := g.calls[key]; ok {
		g.lock.Unlock()
		c.wg.Wait()
		return c.val, false, c.err
	}

	c := g.makeCall(key, fn)
	return c.val, true, c.err
}

func (g *sharedGroup) makeCall(key string, fn func() (interface{}, error)) *call {
	c := new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.lock.Unlock()

	defer func() {
		// delete key first, done later. can't reverse the order, because if reverse,
		// another Do call might wg.Wait() without get notified with wg.Done()
		g.lock.Lock()
		delete(g.calls, key)
		g.lock.Unlock()
		c.wg.Done()
	}()

	c.val, c.err = fn()
	return c
}
