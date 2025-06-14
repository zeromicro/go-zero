package types

import "context"

type (
	Producer interface {
		Send(ctx context.Context, messages ...*Message) error
		SendDelay(ctx context.Context, delaySeconds int64, messages ...*Message) error
		Close() error
	}

	ConsumerGroup interface {
		Start()
		Close() error
		MarkMessage(message *Message)
		Commit()
	}
)
