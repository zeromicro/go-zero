//go:generate mockgen -package mongo -destination iter_mock.go -source iter.go Iter

package mongo

import (
	"github.com/globalsign/mgo/bson"
	"github.com/zeromicro/go-zero/core/breaker"
)

type (
	// Iter interface represents a mongo iter.
	Iter interface {
		All(result interface{}) error
		Close() error
		Done() bool
		Err() error
		For(result interface{}, f func() error) error
		Next(result interface{}) bool
		State() (int64, []bson.Raw)
		Timeout() bool
	}

	// A ClosableIter is a closable mongo iter.
	ClosableIter struct {
		Iter
		Cleanup func()
	}

	promisedIter struct {
		Iter
		promise keepablePromise
	}

	rejectedIter struct{}
)

func (i promisedIter) All(result interface{}) error {
	return i.promise.keep(i.Iter.All(result))
}

func (i promisedIter) Close() error {
	return i.promise.keep(i.Iter.Close())
}

func (i promisedIter) Err() error {
	return i.Iter.Err()
}

func (i promisedIter) For(result interface{}, f func() error) error {
	var ferr error
	err := i.Iter.For(result, func() error {
		ferr = f()
		return ferr
	})
	if ferr == err {
		return i.promise.accept(err)
	}

	return i.promise.keep(err)
}

// Close closes a mongo iter.
func (it *ClosableIter) Close() error {
	err := it.Iter.Close()
	it.Cleanup()
	return err
}

func (i rejectedIter) All(result interface{}) error {
	return breaker.ErrServiceUnavailable
}

func (i rejectedIter) Close() error {
	return breaker.ErrServiceUnavailable
}

func (i rejectedIter) Done() bool {
	return false
}

func (i rejectedIter) Err() error {
	return breaker.ErrServiceUnavailable
}

func (i rejectedIter) For(result interface{}, f func() error) error {
	return breaker.ErrServiceUnavailable
}

func (i rejectedIter) Next(result interface{}) bool {
	return false
}

func (i rejectedIter) State() (int64, []bson.Raw) {
	return 0, nil
}

func (i rejectedIter) Timeout() bool {
	return false
}
