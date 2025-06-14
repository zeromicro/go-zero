package kafka

import (
	"context"
	"log"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logx"
	zsarama "github.com/zeromicro/go-zero/kafka/internal/sarama"
)

const (
	// OffsetNewest stands for the log head offset on the broker for a
	// partition.
	OffsetNewest int64 = -1
	// OffsetOldest stands for the oldest offset available on the broker for a
	// partition.
	OffsetOldest int64 = -2
)

type (
	// Consumer is a wrapper of sarama.Consumer.
	Consumer interface {
		// GetNewestOffset returns the newest offset of the partition.
		GetNewestOffset(topic string, partition int32) (int64, error)
		// ConsumePartition wraps sarama.consumer.ConsumePartition with handlers,
		// errorHandler is optional, if it is nil, the error will be logged.
		ConsumePartition(
			topic string, partition int32, offset int64,
			handler ConsumerHandler, errorHandler ConsumerErrorHandler,
		) (PartitionConsumer, error)
		// Close closes the consumer,
		// the underline sarama.Consumer.Close will be called.
		Close() error
	}
	// PartitionConsumer is a wrapper of sarama.PartitionConsumer.
	PartitionConsumer = zsarama.PartitionConsumer
	// ConsumerHandler is a handler for consuming messages.
	ConsumerHandler func(ctx context.Context, message *Message)
	// ConsumerError is type alias of sarama.ConsumerError.
	ConsumerError = zsarama.ConsumerError
	// ConsumerErrorHandler is a handler for consuming errors.
	ConsumerErrorHandler func(consumerError *ConsumerError)

	consumerWrapper struct {
		zsarama.Consumer
	}
)

// MustNewConsumer is creating a new consumer or die.
func MustNewConsumer(config ConsumerConfig) Consumer {
	c, err := NewConsumer(config)
	if err != nil {
		log.Fatalf("%+v\n\n%s", err, debug.Stack())
	}
	return c
}

// NewConsumer creates a new consumer.
func NewConsumer(config ConsumerConfig) (Consumer, error) {
	logx.Infof("zkafka new consumer, brokers: %v", config.Client.Brokers)

	c, err := zsarama.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	return &consumerWrapper{c}, nil
}

func (c *consumerWrapper) ConsumePartition(topic string, partition int32, offset int64,
	handler ConsumerHandler, errorHandler ConsumerErrorHandler) (PartitionConsumer, error) {
	return c.Consumer.ConsumePartition(topic, partition, offset,
		zsarama.ConsumerHandler(handler), zsarama.ConsumerErrorHandler(errorHandler))
}
