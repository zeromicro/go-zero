package syncx

import "sync"

// Once returns a func that guarantees fn can only called once.
func Once(fn func()) func() {
	once := new(sync.Once)
	return func() {
		once.Do(fn)
	}
}
