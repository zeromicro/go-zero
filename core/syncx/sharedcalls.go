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
		wg       sync.WaitGroup
		consumer sync.WaitGroup
		val      interface{}
		err      error
	}

	sharedGroup struct {
		calls *sync.Map
		pools *sync.Pool //
	}
)

// NewSharedCalls returns a SharedCalls.
func NewSharedCalls() SharedCalls {
	newFun := func() interface{} {
		return &call{}
	}
	return &sharedGroup{
		calls: &sync.Map{},
		pools: &sync.Pool{New: newFun},
	}
}

func (g *sharedGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	val, _, err := g.do(key, fn)
	return val, err
}

func (g *sharedGroup) DoEx(key string, fn func() (interface{}, error)) (val interface{}, fresh bool, err error) {
	return g.do(key, fn)
}
func (g *sharedGroup) do(key string, fn func() (interface{}, error)) (val interface{}, fresh bool, err error) {
	// get a *call
	cc := g.pools.Get()
	if v, ok := g.calls.LoadOrStore(key, cc); ok {
		// restore a *call
		g.pools.Put(cc)

		c := v.(*call)
		c.wg.Wait()

		c.consumer.Add(1)
		val = c.val
		err = c.err
		c.consumer.Done()
		return val, false, err
	} else {
		c := v.(*call)
		defer func() {
			c.wg.Done()
			g.calls.Delete(key)

			// wait for the consumption of other goroutines to complete,restore a *call
			c.consumer.Wait()
			g.pools.Put(c)

		}()
		c.wg.Add(1)
		c.val, c.err = fn()
		return c.val, true, c.err
	}
}
