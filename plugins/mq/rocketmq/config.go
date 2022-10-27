package rocketmq

import (
	"errors"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProducerConf struct {
	NsResolver                 []string // resolver address e.g. 127.0.0.1:9876
	GroupName                  string   `json:",optional"`
	Namespace                  string   `json:",optional"`
	InstanceName               string   `json:",optional"`
	MsgTimeOut                 int      `json:",optional"`
	DefaultTopicQueueNums      int      `json:",optional"`
	CreateTopicKey             string   `json:",optional"`
	CompressMsgBodyOverHowMuch int      `json:",optional"`
	CompressLevel              int      `json:",optional"`
	Retry                      int      `json:",optional"`
	AccessKey                  string   `json:",optional"`
	SecretKey                  string   `json:",optional"`
}

// Validate the configuration
// If the configuration is not set, change it into default
func (c *ProducerConf) Validate() error {
	if c.NsResolver == nil {
		return errors.New("the revolver must not be empty")
	}
	if c.GroupName == "" {
		c.GroupName = "DEFAULT_PRODUCER"
	}
	if c.Namespace == "" {
		c.Namespace = "DEFAULT"
	}
	if c.InstanceName == "" {
		c.InstanceName = "DEFAULT"
	}
	if c.MsgTimeOut == 0 {
		c.MsgTimeOut = 3
	}
	if c.DefaultTopicQueueNums == 0 {
		c.DefaultTopicQueueNums = 4
	}
	if c.CreateTopicKey == "" {
		c.CreateTopicKey = "TBW102"
	}
	if c.CompressMsgBodyOverHowMuch == 0 {
		c.CompressMsgBodyOverHowMuch = 4096
	}
	if c.CompressLevel == 0 {
		c.CompressLevel = 5
	}
	if c.Retry == 0 {
		c.Retry = 2
	}
	return nil
}

func (c *ProducerConf) NewProducer() rocketmq.Producer {
	err := c.Validate()
	logx.Must(err)

	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(c.NsResolver)),
		producer.WithGroupName(c.GroupName),
		producer.WithNamespace(c.Namespace),
		producer.WithInstanceName(c.InstanceName),
		producer.WithSendMsgTimeout(time.Duration(c.MsgTimeOut)*time.Second),
		producer.WithDefaultTopicQueueNums(c.DefaultTopicQueueNums),
		producer.WithCreateTopicKey(c.CreateTopicKey),
		producer.WithCompressMsgBodyOverHowmuch(c.CompressMsgBodyOverHowMuch),
		producer.WithCompressLevel(c.CompressLevel),
		producer.WithRetry(c.Retry),
		producer.WithCredentials(primitive.Credentials{AccessKey: c.AccessKey, SecretKey: c.SecretKey}),
	)

	logx.Must(err)

	return p
}

type ConsumerConf struct {
	NsResolver            []string
	GroupName             string
	Namespace             string
	InstanceName          string
	Strategy              string
	RebalanceLockInterval int
	MaxReconsumeTimes     int32  // 1 means 16 times
	ConsumerModel         string // BroadCasting or Clustering or Unknown
	AutoCommit            bool
	Resolver              string
	AccessKey             string
	SecretKey             string
}

func (c *ConsumerConf) Validate() error {
	if c.NsResolver == nil {
		return errors.New("the revolver must not be empty")
	}
	if c.GroupName == "" {
		c.GroupName = "DEFAULT_CONSUMER"
	}
	if c.Namespace == "" {
		c.Namespace = "DEFAULT"
	}
	if c.InstanceName == "" {
		c.InstanceName = "DEFAULT"
	}
	if c.Strategy == "" {
		c.Strategy = "AllocateByAveragely"
	}
	if c.RebalanceLockInterval == 0 {
		c.RebalanceLockInterval = 20
	}
	if c.MaxReconsumeTimes == 0 {
		c.MaxReconsumeTimes = -1
	}
	if c.ConsumerModel == "" {
		c.ConsumerModel = "Clustering"
	}
	if c.Resolver == "" {
		c.Resolver = "DEFAULT"
	}
	return nil
}

func (c *ConsumerConf) NewPushConsumer() rocketmq.PushConsumer {
	err := c.Validate()
	logx.Must(err)

	var strategy consumer.AllocateStrategy
	switch c.Strategy {
	case "AllocateByAveragely":
		strategy = consumer.AllocateByAveragely
	case "AllocateByAveragelyCircle":
		strategy = consumer.AllocateByAveragelyCircle
	case "AllocateByMachineNearby":
		strategy = consumer.AllocateByMachineNearby
	default:
		strategy = consumer.AllocateByAveragely
	}

	csm, err := rocketmq.NewPushConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver(c.NsResolver)),
		consumer.WithGroupName(c.GroupName),
		consumer.WithNamespace(c.Namespace),
		consumer.WithStrategy(strategy),
		consumer.WithRebalanceLockInterval(time.Duration(c.RebalanceLockInterval)*time.Second),
		consumer.WithMaxReconsumeTimes(c.MaxReconsumeTimes),
		consumer.WithCredentials(primitive.Credentials{AccessKey: c.AccessKey, SecretKey: c.SecretKey}),
		consumer.WithAutoCommit(c.AutoCommit),
		consumer.WithInstance(c.InstanceName),
	)

	logx.Must(err)

	return csm
}

func (c *ConsumerConf) NewPullConsumer() rocketmq.PullConsumer {
	err := c.Validate()
	logx.Must(err)

	var strategy consumer.AllocateStrategy
	switch c.Strategy {
	case "AllocateByAveragely":
		strategy = consumer.AllocateByAveragely
	case "AllocateByAveragelyCircle":
		strategy = consumer.AllocateByAveragelyCircle
	case "AllocateByMachineNearby":
		strategy = consumer.AllocateByMachineNearby
	default:
		strategy = consumer.AllocateByAveragely
	}

	csm, err := rocketmq.NewPullConsumer(
		consumer.WithNsResolver(primitive.NewPassthroughResolver(c.NsResolver)),
		consumer.WithGroupName(c.GroupName),
		consumer.WithNamespace(c.Namespace),
		consumer.WithStrategy(strategy),
		consumer.WithRebalanceLockInterval(time.Duration(c.RebalanceLockInterval)*time.Second),
		consumer.WithMaxReconsumeTimes(c.MaxReconsumeTimes),
		consumer.WithCredentials(primitive.Credentials{AccessKey: c.AccessKey, SecretKey: c.SecretKey}),
		consumer.WithAutoCommit(c.AutoCommit),
		consumer.WithInstance(c.InstanceName),
	)

	logx.Must(err)

	return csm
}
