package syncx

import "sync"

type (
	// SingleFlight lets the concurrent calls with the same key to share the call result.
	// For example, A called F, before it's done, B called F. Then B would not execute F,
	// and shared the result returned by F which called by A.
	// The calls with the same key are dependent, concurrent calls share the returned values.
	// A ------->calls F with key<------------------->returns val
	// B --------------------->calls F with key------>returns val
	SingleFlight interface {
		Do(key string, fn func() (interface{}, error)) (interface{}, error)
		DoEx(key string, fn func() (interface{}, error)) (interface{}, bool, error)
	}

	flightGroup struct {
		calls sync.Map
	}
)

// NewSingleFlight returns a SingleFlight.
func NewSingleFlight() SingleFlight {
	return &flightGroup{
		calls: sync.Map{},
	}
}

func (g *flightGroup) Do(key string, fn func() (interface{}, error)) (val interface{}, err error) {
	val, _, err = g.createCall(key, fn)
	return
}

func (g *flightGroup) DoEx(key string, fn func() (interface{}, error)) (val interface{}, fresh bool, err error) {
	val, fresh, err = g.createCall(key, fn)
	return
}

func (g *flightGroup) createCall(key string, fn func() (interface{}, error)) (val interface{}, fresh bool, err error) {
	val, err = fn()
	if err != nil {
		return
	}
	val, fresh = g.calls.LoadOrStore(key, val)
	return
}

//说明：sync.Map的原子操作，适用于读多很少，相比以前读写都使用mutex，性能会更好
//测试用例部分过不了，所有操作需要先执行fn() 函数，，Kevin老师是否考虑使用sync.Map