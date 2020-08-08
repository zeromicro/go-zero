package syncx

import (
	"errors"

	"github.com/tal-tech/go-zero/core/lang"
)

var ErrReturn = errors.New("discarding limited token, resource pool is full, someone returned multiple times")

type Limit struct {
	pool chan lang.PlaceholderType
}

func NewLimit(n int) Limit {
	return Limit{
		pool: make(chan lang.PlaceholderType, n),
	}
}

func (l Limit) Borrow() {
	l.pool <- lang.Placeholder
}

// Return returns the borrowed resource, returns error only if returned more than borrowed.
func (l Limit) Return() error {
	select {
	case <-l.pool:
		return nil
	default:
		return ErrReturn
	}
}

func (l Limit) TryBorrow() bool {
	select {
	case l.pool <- lang.Placeholder:
		return true
	default:
		return false
	}
}
