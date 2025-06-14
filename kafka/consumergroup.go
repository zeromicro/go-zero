package kafka

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/internal/health"
	zsarama "github.com/zeromicro/go-zero/kafka/internal/sarama"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

const probeNamePrefix = "kafka-consumergroup"

type (
	Handler func(ctx context.Context, message *Message) error

	ConsumerGroup interface {
		service.Service
		MarkMessage(message *Message)
		Commit()
	}

	consumerGroup struct {
		c             types.ConsumerGroup
		healthManager health.Probe
	}
)

// MustNewConsumerGroup returns a ConsumerGroup, exits on any error.
func MustNewConsumerGroup(c ConsumerGroupConfig, handler Handler) ConsumerGroup {
	cg, err := NewConsumerGroup(c, handler)
	if err != nil {
		log.Fatalf("%+v\n\n%s", err, debug.Stack())
	}
	return cg
}

// NewConsumerGroup create kafka consumer group from config.
func NewConsumerGroup(c ConsumerGroupConfig, handler Handler) (ConsumerGroup, error) {
	logx.Infof("zkafka new consumer group, topic: %s, groupId: %s, brokers: %v",
		c.Topic, c.GroupID, c.Client.Brokers)

	sc, err := zsarama.NewConsumerGroup(c, zsarama.Handler(handler))
	if err != nil {
		return nil, err
	}

	return &consumerGroup{
		c:             sc,
		healthManager: health.NewHealthManager(fmt.Sprintf("%s-%s", probeNamePrefix, c.Topic)),
	}, nil
}

func (c *consumerGroup) Start() {
	waitForCalled := proc.AddShutdownListener(func() {
		c.Stop()
	})
	defer waitForCalled()

	c.healthManager.MarkReady()
	// add component probe to global health manager.
	health.AddProbe(c.healthManager)
	c.c.Start()
}

func (c *consumerGroup) Stop() {
	c.healthManager.MarkNotReady()
	if err := c.c.Close(); err != nil {
		logx.Errorf("stop consumer group err: %s", err)
	}
}

// MarkMessage marks a message as consumed, you need call Commit() to commit the offset to the backend
func (c *consumerGroup) MarkMessage(message *Message) {
	c.c.MarkMessage(message)
}

// Commit Note: do not call Commit() in Handler!!! calling Commit() performs a blocking synchronous operation
func (c *consumerGroup) Commit() {
	c.c.Commit()
}
