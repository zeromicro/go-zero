package kafka

import (
	"context"
	"log"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logx"
	zsarama "github.com/zeromicro/go-zero/kafka/internal/sarama"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

type (
	Producer interface {
		Send(ctx context.Context, messages ...*Message) error
		SendDelay(ctx context.Context, delaySeconds int64, messages ...*Message) error
	}

	producer struct {
		p  types.Producer
		pc ProducerConfig
	}
)

// MustNewProducer returns a Producer, exits on any error.
func MustNewProducer(pc ProducerConfig) Producer {
	p, err := NewProducer(pc)
	if err != nil {
		log.Fatalf("%+v\n\n%s", err, debug.Stack())
	}

	return p
}

// NewProducer create kafka Producer from config.
func NewProducer(pc ProducerConfig) (Producer, error) {
	logx.Infof("zkafka new producer, topic: %s, brokers: %v", pc.Topic, pc.Client.Brokers)

	p, err := zsarama.NewProducer(pc)
	if err != nil {
		return nil, err
	}

	return &producer{p: p, pc: pc}, nil
}

// Send is low level API for generated code.
func (c *producer) Send(ctx context.Context, messages ...*Message) error {
	return c.p.Send(ctx, messages...)
}

// SendDelay sends a delay messages to kafka delay queue
func (c *producer) SendDelay(ctx context.Context, delaySeconds int64, messages ...*Message) error {
	return c.p.SendDelay(ctx, delaySeconds, messages...)
}
