package syncx

import "sync"

// Once returns a func that guarantees fn can only called once.
// Deprecated: use sync.OnceFunc instead.
func Once(fn func()) func() {
	return sync.OnceFunc(fn)
}
