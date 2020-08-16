package queue

import (
	"errors"
	"sync/atomic"

	"github.com/tal-tech/go-zero/core/logx"
)

var ErrNoAvailablePusher = errors.New("no available pusher")

type BalancedPusher struct {
	name    string
	pushers []Pusher
	index   uint64
}

func NewBalancedPusher(pushers []Pusher) Pusher {
	return &BalancedPusher{
		name:    generateName(pushers),
		pushers: pushers,
	}
}

func (pusher *BalancedPusher) Name() string {
	return pusher.name
}

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
