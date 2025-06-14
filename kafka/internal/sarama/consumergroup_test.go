package sarama

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/retry"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

func Test_fillConsumerConfig(t *testing.T) {
	scc := types.SharedConsumerConfig{
		BalanceStrategy:   sarama.RangeBalanceStrategyName,
		MaxProcessingTime: 2 * time.Second,
	}
	sc := &sarama.Config{}
	err := fillConsumerConfig(sc, scc)
	assert.NoError(t, err, fmt.Sprintf("toSaramaConfig"))
	assert.Equal(t, 2*time.Second, sc.Consumer.MaxProcessingTime)
}

func TestNewConsumerGroup(t *testing.T) {
	t.Run("invalid params", func(t *testing.T) {
		_, err := NewConsumerGroup(types.ConsumerGroupConfig{}, nil)
		assert.Error(t, err)

		_, err = NewConsumerGroup(types.ConsumerGroupConfig{
			Client:      types.ClientConfig{Brokers: []string{"b1"}},
			GroupConfig: types.GroupConfig{InitialOffset: "xx"},
		}, nil)
		assert.Error(t, err)
	})

	t.Run("new client fail", func(t *testing.T) {
		_, err := NewConsumerGroup(types.ConsumerGroupConfig{
			Client: types.ClientConfig{Brokers: []string{"b1"}},
		}, nil)
		assert.Error(t, err)
	})
}

func TestClientCancel(t *testing.T) {
	runWithClient(t, func(t *testing.T, client *Client) {
		t.Run("kafka_consumer_group cancel is nil", func(t *testing.T) {
			topic := "kafka_consumer_group"
			group, err := client.NewConsumerGroup(types.GroupConfig{
				Topic: topic,
			}, nil)
			assert.Nil(t, group.cancel)
			assert.Nil(t, err)
		})

		t.Run("kafka_consumer_group cancel exist", func(t *testing.T) {
			topic := "kafka_consumer_group"
			group, _ := client.NewConsumerGroup(types.GroupConfig{
				Topic: topic,
			}, func(ctx context.Context, message *types.Message) error {
				return nil
			})
			group.cancel = func() {}
			group.Close()
		})
	})
}

func TestConsumeTTT(t *testing.T) {
	runConsumerGroupConsume(true, true, true, t)
	//runConsumerGroupConsume(false, true, true, t)
	//runConsumerGroupConsume(false, false, true, t)
	//runConsumerGroupConsume(false, false, false, t)
}

func runConsumerGroupConsume(ConsumeTimeout, EnableRecovery, AutoCommit bool, t *testing.T) {
	//logx.Disable()
	var c trace.Config
	conf.FillDefault(&c)
	trace.StartAgent(c)

	broker0 := sarama.NewMockBroker(t, 1)
	broker0.SetHandlerByMap(map[string]sarama.MockResponse{
		"ApiVersionsRequest": sarama.NewMockApiVersionsResponse(t),
		"MetadataRequest": sarama.NewMockMetadataResponse(t).
			SetBroker(broker0.Addr(), broker0.BrokerID()).
			SetController(broker0.BrokerID()).
			SetLeader("my-topic", 0, broker0.BrokerID()),
		"OffsetRequest": sarama.NewMockOffsetResponse(t).
			SetOffset("my-topic", 0, sarama.OffsetOldest, 0).
			SetOffset("my-topic", 0, sarama.OffsetNewest, 1),
		"FindCoordinatorRequest": sarama.NewMockFindCoordinatorResponse(t).
			SetCoordinator(sarama.CoordinatorGroup, "my-group", broker0),
		"HeartbeatRequest": sarama.NewMockHeartbeatResponse(t),
		"JoinGroupRequest": sarama.NewMockSequence(
			sarama.NewMockJoinGroupResponse(t).SetError(sarama.ErrOffsetsLoadInProgress),
			sarama.NewMockJoinGroupResponse(t).SetGroupProtocol(sarama.RangeBalanceStrategyName),
		),
		"SyncGroupRequest": sarama.NewMockSequence(
			sarama.NewMockSyncGroupResponse(t).SetError(sarama.ErrOffsetsLoadInProgress),
			sarama.NewMockSyncGroupResponse(t).SetMemberAssignment(
				&sarama.ConsumerGroupMemberAssignment{
					Version: 0,
					Topics: map[string][]int32{
						"my-topic": {0},
					},
				}),
		),
		"OffsetFetchRequest": sarama.NewMockOffsetFetchResponse(t).SetOffset(
			"my-group", "my-topic", 0, 0, "", sarama.ErrNoError,
		).SetError(sarama.ErrNoError),
		"FetchRequest": sarama.NewMockSequence(
			sarama.NewMockFetchResponse(t, 1).
				SetMessage("my-topic", 0, 0, sarama.StringEncoder("foo")).
				SetMessage("my-topic", 0, 1, sarama.StringEncoder("bar")),
			sarama.NewMockFetchResponse(t, 1),
		),
	})
	defer broker0.Close()

	var config types.UniversalClientConfig
	_ = conf.FillDefault(&config)
	config.Client.Brokers = []string{broker0.Addr()}
	config.Client.ClientId = t.Name()
	if ConsumeTimeout {
		config.Consumer.ConsumeTimeout = time.Second
	}

	var (
		err         error
		groupConfig types.GroupConfig
	)
	_ = conf.FillDefault(&groupConfig)
	groupConfig.GroupID = "my-group"
	groupConfig.Topic = "my-topic"
	groupConfig.EnableRecovery = EnableRecovery
	groupConfig.AutoCommit = AutoCommit

	sc, err := toSaramaConfig(config.Client)
	assert.Nil(t, err)
	sc.Version = sarama.V2_0_0_0

	err = fillProducerConfig(sc, config.Producer)
	assert.Nil(t, err)
	err = fillConsumerConfig(sc, config.Consumer)
	assert.Nil(t, err)

	cc, err := sarama.NewClient(config.Client.Brokers, sc)
	assert.Nil(t, err)

	saramaClient := &Client{
		saramaClient: cc,
		config:       config,
		saramaConfig: sc,
	}
	assert.Nil(t, err)
	defer saramaClient.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	var s *ConsumerGroup

	rc := retry.DefaultConfig()
	rc.MaxRetries = groupConfig.RetryConfig.MaxRetries
	sc.Consumer.Return.Errors = false // consumergroup不处理 sarama -> Errors() <-chan error. 默认println err.

	initialOffset, err := parseInitialOffset(groupConfig.InitialOffset)
	assert.Nil(t, err)
	sc.Consumer.Offsets.Initial = initialOffset

	s = &ConsumerGroup{
		name:                 config.Client.GetClientName(),
		clientConsumerConfig: config.Consumer,
		consumerGroupConfig:  groupConfig,
		consumeRetryInterval: time.Second,
		backOffConfig:        rc,
		handler: func(ctx context.Context, message *types.Message) error {
			fmt.Println("received a message")
			logc.Info(ctx, logx.Field("message", message))
			s.MarkMessage(message)
			s.Commit()
			wg.Done()
			return nil
		},
		sc: *sc,
	}

	err = s.initClient(*sc, []string{broker0.Addr()})
	assert.Nil(t, err)
	defer func() { _ = s.Close() }()
	go func() {
		s.Start()
	}()
	wg.Wait()
	assert.Nil(t, err)
}
