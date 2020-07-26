package queue

import (
	"errors"
	"sync/atomic"

	"zero/core/logx"
)

var ErrNoAvailablePusher = errors.New("no available pusher")

type BalancedQueuePusher struct {
	name    string
	pushers []QueuePusher
	index   uint64
}

func NewBalancedQueuePusher(pushers []QueuePusher) QueuePusher {
	return &BalancedQueuePusher{
		name:    generateName(pushers),
		pushers: pushers,
	}
}

func (pusher *BalancedQueuePusher) Name() string {
	return pusher.name
}

func (pusher *BalancedQueuePusher) Push(message string) error {
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
