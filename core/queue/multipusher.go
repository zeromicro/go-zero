package queue

import "github.com/zeromicro/go-zero/core/errorx"

// A MultiPusher is a pusher that can push messages to multiple underlying pushers.
type MultiPusher struct {
	name    string
	pushers []Pusher
}

// NewMultiPusher returns a MultiPusher.
func NewMultiPusher(pushers []Pusher) Pusher {
	return &MultiPusher{
		name:    generateName(pushers),
		pushers: pushers,
	}
}

// Name returns the name of pusher.
func (pusher *MultiPusher) Name() string {
	return pusher.name
}

// Push pushes a message into the underlying queue.
func (pusher *MultiPusher) Push(message string) error {
	var batchError errorx.BatchError

	for _, each := range pusher.pushers {
		if err := each.Push(message); err != nil {
			batchError.Add(err)
		}
	}

	return batchError.Err()
}
