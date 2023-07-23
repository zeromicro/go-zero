package queue

import (
	"errors"
	"sync/atomic"

	"github.com/zeromicro/go-zero/core/logx"
)

// ErrNoAvailablePusher indicates no pusher available.
var ErrNoAvailablePusher = errors.New("no available pusher")

// A BalancedPusher is used to push messages to multiple pusher with round robin algorithm.
type BalancedPusher struct {
	name    string
	pushers []Pusher
	index   uint64
}

// NewBalancedPusher returns a BalancedPusher.
func NewBalancedPusher(pushers []Pusher) Pusher {
	return &BalancedPusher{
		name:    generateName(pushers),
		pushers: pushers,
	}
}

// Name returns the name of pusher.
func (pusher *BalancedPusher) Name() string {
	return pusher.name
}

// Push pushes message to one of the underlying pushers.
func (pusher *BalancedPusher) Push(message string) error {
	size := len(pusher.pushers)

	for i := 0; i < size; i++ {
		index := atomic.AddUint64(&pusher.index, 1) % uint64(size)
		target := pusher.pushers[index]

		if err := target.Push(message); err != nil {
			logx.Error(err)
		} else {
			return nil
		}
	}

	return ErrNoAvailablePusher
}
