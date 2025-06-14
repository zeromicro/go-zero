package sarama

import (
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

type Client struct {
	saramaClient              sarama.Client
	saramaConfig              *sarama.Config
	config                    types.UniversalClientConfig
	isExiting                 bool
	pwg                       sync.WaitGroup
	cwg                       sync.WaitGroup
	consumers                 sync.Map
	consumerGroups            sync.Map
	producers                 sync.Map
	k                         string
	spawnedPartitionConsumers sync.Map // [string]UnsafePartitionConsumer
}

func NewClient(c types.UniversalClientConfig) (*Client, error) {
	sc, err := toSaramaConfig(c.Client)
	if err != nil {
		return nil, err
	}

	if err := fillProducerConfig(sc, c.Producer); err != nil {
		return nil, err
	}
	if err := fillConsumerConfig(sc, c.Consumer); err != nil {
		return nil, err
	}

	cc, err := sarama.NewClient(c.Client.Brokers, sc)
	if err != nil {
		return nil, err
	}

	return &Client{
		saramaClient: cc,
		config:       c,
		saramaConfig: sc,
	}, nil
}

func (c *Client) NewProducer(defaultTopic string) (*Producer, error) {
	return newProducerFromClient(c, defaultTopic)
}

func (c *Client) NewAsyncProducer(defaultTopic string, opt AsyncProducerOption) (*AsyncProducer, error) {
	cc, err := newAsyncProducerFromClient(c, defaultTopic, opt)
	if err != nil {
		return nil, err
	}
	cc.setupCallbacks()
	return cc, nil
}

func (c *Client) NewConsumerGroup(cc types.GroupConfig, handler Handler) (
	*ConsumerGroup, error) {
	cg, err := NewConsumerGroup(types.ConsumerGroupConfig{
		Client:               c.config.Client,
		SharedConsumerConfig: c.config.Consumer,
		GroupConfig:          cc,
	}, handler)
	if err == nil {
		c.consumerGroups.Store(cg.name, cg)
	}
	return cg, err
}

func (c *Client) NewConsumer() (Consumer, error) {
	return newConsumerFromClient(c)
}

func (c *Client) waitTimeout(t time.Duration, waitFuc func(), wg *sync.WaitGroup) {
	var ch = make(chan struct{})
	go func() {
		waitFuc()
	}()
	go func() {
		wg.Wait() // 等待所有的wg结束
		ch <- struct{}{}
	}()
	select {
	case <-time.After(t):
		return
	case <-ch:
		return
	}
}

func (c *Client) Close() error {
	total := 5000 * time.Millisecond
	begin := time.Now()

	c.isExiting = true
	c.waitTimeout(total, func() {
		c.consumers.Range(func(key, value any) bool {
			c.consumers.Delete(key)
			_ = value.(UnsafeConsumer).Close()
			return true
		})
		c.consumerGroups.Range(func(key, value any) bool {
			c.consumerGroups.Delete(key)
			_ = value.(*ConsumerGroup).Close()
			return true
		})
	}, &c.cwg) // 等待所有的consumer结束

	c.waitTimeout(total-time.Since(begin), func() {
		c.producers.Range(func(key, value any) bool {
			c.producers.Delete(key)
			_ = value.(*UnsafeProducer).Close()
			return true
		})
	}, &c.pwg) // 等待所有的producer结束
	return c.saramaClient.Close()
}

func (c *Client) GetOffset(topic string, partition int32, time int64) (int64, error) {
	return c.saramaClient.GetOffset(topic, partition, time)
}

func (c *Client) Partitions(topic string) ([]int32, error) {
	return c.saramaClient.Partitions(topic)
}
