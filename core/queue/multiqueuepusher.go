package queue

import "github.com/tal-tech/go-zero/core/errorx"

type MultiQueuePusher struct {
	name    string
	pushers []Pusher
}

func NewMultiQueuePusher(pushers []Pusher) Pusher {
	return &MultiQueuePusher{
		name:    generateName(pushers),
		pushers: pushers,
	}
}

func (pusher *MultiQueuePusher) Name() string {
	return pusher.name
}

func (pusher *MultiQueuePusher) Push(message string) error {
	var batchError errorx.BatchError

	for _, each := range pusher.pushers {
		if err := each.Push(message); err != nil {
			batchError.Add(err)
		}
	}

	return batchError.Err()
}
