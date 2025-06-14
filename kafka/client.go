package kafka

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/zeromicro/go-zero/internal/health"
	zsarama "github.com/zeromicro/go-zero/kafka/internal/sarama"
)

type (
	Client interface {
		NewProducer() (Producer, error)
		NewAsyncProducer() (AsyncProducer, error)
		NewConsumerGroup(cc GroupConfig, handler Handler) (ConsumerGroup, error)
		NewConsumer() (Consumer, error)
		// GetOffset returns the offset of the partition.
		GetOffset(topic string, partition int32, time int64) (int64, error)
		// Partitions returns the partitions of the topic.
		Partitions(topic string) ([]int32, error)
		Close() error
	}

	client struct {
		*zsarama.Client
	}
)

func MustNewClient(config UniversalClientConfig) Client {
	c, err := NewClient(config)
	if err != nil {
		log.Fatalf("%+v\n\n%s", err, debug.Stack())
	}
	return c
}

func NewClient(config UniversalClientConfig) (Client, error) {
	c, err := zsarama.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

func (c *client) NewProducer() (Producer, error) {
	return c.Client.NewProducer("")
}

func (c *client) NewAsyncProducer() (AsyncProducer, error) {
	return c.Client.NewAsyncProducer("", zsarama.AsyncProducerOption{})
}

func (c *client) NewConsumerGroup(cc GroupConfig, handler Handler) (ConsumerGroup, error) {
	sc, err := c.Client.NewConsumerGroup(cc, zsarama.Handler(handler))
	if err != nil {
		return nil, err
	}

	return &consumerGroup{
		c:             sc,
		healthManager: health.NewHealthManager(fmt.Sprintf("%s-%s", probeNamePrefix, cc.Topic)),
	}, nil
}

func (c *client) NewConsumer() (Consumer, error) {
	cc, err := c.Client.NewConsumer()
	if err != nil {
		return nil, err
	}
	return &consumerWrapper{cc}, nil
}
