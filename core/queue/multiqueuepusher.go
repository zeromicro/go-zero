package queue

import "zero/core/errorx"

type MultiQueuePusher struct {
	name    string
	pushers []QueuePusher
}

func NewMultiQueuePusher(pushers []QueuePusher) QueuePusher {
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
