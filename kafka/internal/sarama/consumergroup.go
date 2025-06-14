package sarama

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/cenkalti/backoff/v4"
	"github.com/dnwe/otelsarama"
	"github.com/zeromicro/go-zero/core/contextx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/retry"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/kafka/internal/metrics"
	ztrace "github.com/zeromicro/go-zero/kafka/internal/trace"
	"github.com/zeromicro/go-zero/kafka/internal/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otrace "go.opentelemetry.io/otel/trace"
)

var _ types.ConsumerGroup = (*ConsumerGroup)(nil)
var lessLogger = logx.NewLessLogger(1000)

type (
	Handler func(ctx context.Context, message *types.Message) error

	ConsumerGroup struct {
		name                 string
		clientConsumerConfig types.SharedConsumerConfig
		consumerGroupConfig  types.GroupConfig
		handler              Handler
		backOffConfig        retry.Config
		consumeRetryInterval time.Duration
		cancel               context.CancelFunc
		cg                   sarama.ConsumerGroup
		sc                   sarama.Config
		client               sarama.Client
		once                 sync.Once
		session              sarama.ConsumerGroupSession
	}

	nextPartitionOffset struct {
		Topic      string
		Partition  int32
		NextOffset int64
		Message    string
	}
)

func NewConsumerGroup(consumerGroupConfig types.ConsumerGroupConfig, handler Handler) (
	*ConsumerGroup, error) {
	rc := retry.DefaultConfig()
	rc.MaxRetries = consumerGroupConfig.RetryConfig.MaxRetries

	sc, err := toSaramaConfig(consumerGroupConfig.Client)
	if err != nil {
		return nil, err
	}
	sc.Consumer.Return.Errors = false // consumergroup不处理 sarama -> Errors() <-chan error. 默认println err.

	initialOffset, err := parseInitialOffset(consumerGroupConfig.InitialOffset)
	if err != nil {
		return nil, err
	}
	sc.Consumer.Offsets.Initial = initialOffset
	sc.Consumer.Offsets.AutoCommit.Enable = consumerGroupConfig.AutoCommit
	if !consumerGroupConfig.AutoCommit {
		consumerGroupConfig.LogOffsets = false
	}

	s := &ConsumerGroup{
		name:                 consumerGroupConfig.Client.GetClientName(),
		clientConsumerConfig: consumerGroupConfig.SharedConsumerConfig,
		consumerGroupConfig:  consumerGroupConfig.GroupConfig,
		consumeRetryInterval: time.Second,
		backOffConfig:        rc,
		handler:              handler,
		sc:                   *sc,
	}

	if err := s.initClient(*sc, consumerGroupConfig.Client.Brokers); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *ConsumerGroup) Start() {
	if s.handler == nil {
		logx.Error("no handler")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	topic := s.consumerGroupConfig.Topic

	logx.Infow(fmt.Sprintf("Subscribed and listening to topic: %v", topic))
	for {
		// If the context was canceled, as is the case when handling SIGINT and SIGTERM below,
		//then this pops us out of the consumer loop.
		if ctx.Err() != nil {
			break
		}

		logx.Infow("Starting loop to consume.")

		// Consume the requested topics
		bo := backoff.WithContext(backoff.NewConstantBackOff(s.consumeRetryInterval), ctx)

		innerErr := retry.NotifyRecover(func() error {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return backoff.Permanent(ctxErr)
			}
			e := s.cg.Consume(ctx, []string{topic}, s)
			if e != nil {
				lessLogger.Errorf("Error consuming %s, internal retry, err:%v", topic, e)
			}
			return e
		}, bo, func(err error, t time.Duration, times int) {
			logx.Errorw(fmt.Sprintf("Error consuming %s. Retrying...: %v", topic, err),
				logx.Field("brokers", s.name))
		}, func(times int) {
			logx.Infof("Recovered consuming %s", topic)
		}, false)

		if innerErr != nil && !errors.Is(innerErr, context.Canceled) {
			logx.Errorw(fmt.Sprintf("Permanent error consuming %s, error: %s",
				s.consumerGroupConfig.Topic, innerErr),
				logx.Field("brokers", s.name))
		}
	}
}

func (s *ConsumerGroup) Close() error {
	if s.cg == nil {
		return nil
	}

	var err error
	s.once.Do(func() {
		logx.Info("consumer group close")
		if s.cancel != nil {
			s.cancel()
		}
		// close ConsumerGroup is not necessary when using sarama.NewConsumerGroupFromClient.
		err = s.client.Close()
	})

	return err
}

func (s *ConsumerGroup) Setup(session sarama.ConsumerGroupSession) error {
	s.session = session
	if !s.consumerGroupConfig.LogOffsets {
		return nil
	}

	nextPartitionOffsets := s.getPartitionOffsets(session)
	logx.Infof("consumer group setup, name=%s, topic=%s, Claims=%v, MemberID=%s, GenerationID=%d, nextPartitionOffsets=%+v", s.name, s.consumerGroupConfig.Topic, session.Claims(), session.MemberID(), session.GenerationID(), nextPartitionOffsets)
	return nil
}

func (s *ConsumerGroup) Cleanup(session sarama.ConsumerGroupSession) error {
	if !s.consumerGroupConfig.LogOffsets {
		return nil
	}

	nextPartitionOffsets := s.getPartitionOffsets(session)
	logx.Infof("consumer group cleanup, name=%s, topic=%s, Claims=%v, MemberID=%s, GenerationID=%d, nextPartitionOffsets=%+v", s.name, s.consumerGroupConfig.Topic, session.Claims(), session.MemberID(), session.GenerationID(), nextPartitionOffsets)
	return nil
}

func (s *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	b := s.backOffConfig.NewBackOffWithContext(session.Context())

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			start := time.Now()

			// extract span info from consumer message into ctx.
			ctx := otel.GetTextMapPropagator().Extract(session.Context(), otelsarama.NewConsumerMessageCarrier(message))
			// remove session.Context() cancel
			ctx = contextx.ValueOnlyFrom(ctx)
			ctx = logx.ContextWithFields(ctx, logx.Field("brokers", s.name))
			logger := logx.WithContext(ctx)

			trace.DoWithSpan(ctx, "kafka.consumergroup.handler", func(ctx context.Context, span otrace.Span) {
				msg := consumerMessage2Message(message)
				ztrace.AddConsumerAttrs(span, msg, attribute.KeyValue{
					Key:   types.KafkaBroker,
					Value: attribute.StringValue(s.name),
				})
				messageKey := string(message.Key)

				err := retry.NotifyRecover(func() error {
					if !s.consumerGroupConfig.EnableRecovery {
						return s.handlerWrapper(ctx, msg)
					}
					return threading.RunSafeWrap(func() error {
						return s.handlerWrapper(ctx, msg)
					})
				}, b, func(err error, duration time.Duration, times int) {
					logger.Errorf("Error processing Kafka message: %s/%d/%d [key=%s]. Error: %v. Retrying...",
						message.Topic, message.Partition, message.Offset, messageKey, err)
				}, func(times int) {
					logger.Infof("Successfully processed Kafka message after it previously failed: %s/%d/%d [key=%s]",
						message.Topic, message.Partition, message.Offset, messageKey)
				}, false)

				cost := time.Since(start)
				metrics.KafkaConsumerGroupHistogram.ObserveFloat(cost.Seconds(), s.name,
					s.consumerGroupConfig.GroupID, message.Topic)
				metrics.KafkaConsumerGroupGauge.Set(cost.Seconds(), metrics.ConsumerGroupKey{
					Brokers: s.name,
					GroupId: s.consumerGroupConfig.GroupID,
					Topic:   message.Topic,
				})

				if err != nil {
					span.RecordError(err)
					metrics.KafkaConsumerGroupCounter.Inc(s.name, s.consumerGroupConfig.GroupID,
						message.Topic, metrics.CodeError)
					logger.Errorw(fmt.Sprintf("kafka handler message error: topic=%s partition=%d offset=%d key=%s error=%v",
						message.Topic, message.Partition, message.Offset, messageKey, err),
						logx.Field("message", string(msg.Value)),
					)
				} else {
					metrics.KafkaConsumerGroupCounter.Inc(s.name, s.consumerGroupConfig.GroupID,
						message.Topic, metrics.CodeOK)
				}
				if !s.consumerGroupConfig.DisableAutoMark {
					session.MarkMessage(message, "")
				}
			}, otrace.WithSpanKind(otrace.SpanKindConsumer))
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			logx.Infof("sarama session done, name: %s, topic: %s, ctx error: %v", s.name,
				s.consumerGroupConfig.Topic, session.Context().Err())
			return nil
		}
	}
}

func (s *ConsumerGroup) handlerWrapper(ctx context.Context, msg *types.Message) error {
	// inject timeout to handler context if needed
	if s.clientConsumerConfig.ConsumeTimeout > 0 {
		handlerCtx, cancel := context.WithTimeout(ctx, s.clientConsumerConfig.ConsumeTimeout)
		defer cancel()
		return s.handler(handlerCtx, msg)
	}
	return s.handler(ctx, msg)
}

func (s *ConsumerGroup) MarkMessage(message *types.Message) {
	if !s.consumerGroupConfig.DisableAutoMark {
		logx.Errorf("MarkMessage is unnecessary when AutoMark set true")
		return
	}
	s.session.MarkOffset(message.Topic, int32(message.Partition), message.Offset+1, "")
}

func (s *ConsumerGroup) Commit() {
	if s.consumerGroupConfig.AutoCommit {
		logx.Errorf("Commit is unnecessary when AutoCommit set true")
		return
	}
	s.session.Commit()
}

func (s *ConsumerGroup) initClient(sc sarama.Config, brokers []string) error { //nolint
	if err := fillConsumerConfig(&sc, s.clientConsumerConfig); err != nil {
		return err
	}

	// PLEASE NOTE: consumer groups can only re-use but not share clients.
	client, err := sarama.NewClient(brokers, &sc)
	if err != nil {
		return err
	}
	s.client = client

	return s.initConsumerGroup()
}

func (s *ConsumerGroup) initConsumerGroup() error {
	cg, err := sarama.NewConsumerGroupFromClient(s.consumerGroupConfig.GroupID, s.client)
	if err != nil {
		return err
	}
	s.cg = cg
	return nil
}

func (s *ConsumerGroup) getPartitionOffsets(session sarama.ConsumerGroupSession) []nextPartitionOffset {
	manager, err := sarama.NewOffsetManagerFromClient(s.consumerGroupConfig.GroupID, s.client)
	if err != nil {
		logx.Errorf("sarama NewOffsetManagerFromClient error: %v", err)
		return nil
	}
	defer manager.Close()

	var nextPartitionOffsets []nextPartitionOffset
	var lock sync.Mutex

	var wg sync.WaitGroup
	for topic, partitions := range session.Claims() {
		topic := topic
		for _, p := range partitions {
			wg.Add(1)
			partition := p

			go func() {
				defer wg.Done()
				pom, err := manager.ManagePartition(topic, partition)
				if err != nil {
					logx.Errorf("sarama ManagePartition(%s, %d) error: %v", topic, partition, err)
					return
				}

				next, msg := pom.NextOffset()
				lock.Lock()
				nextPartitionOffsets = append(nextPartitionOffsets, nextPartitionOffset{
					Topic:      topic,
					Partition:  partition,
					NextOffset: next,
					Message:    msg,
				})
				lock.Unlock()

				_ = pom.Close()
			}()
		}
	}
	wg.Wait()

	return nextPartitionOffsets
}

func fillConsumerConfig(sc *sarama.Config, cc types.SharedConsumerConfig) error {
	balanceStrategy, err := parseBalanceStrategy(cc.BalanceStrategy)
	if err != nil {
		return err
	}

	// consumer.Group.Rebalance.Strategy is Deprecated, we use GroupStrategies now.
	sc.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{balanceStrategy}
	sc.Consumer.Fetch.Min = int32(cc.FetchMinBytes)      // default is 1
	sc.Consumer.Fetch.Max = int32(cc.FetchMaxBytes)      // default is 10MB
	sc.Consumer.MaxWaitTime = cc.FetchMaxWaitTime        // default is 500ms
	sc.Consumer.MaxProcessingTime = cc.MaxProcessingTime //default is 100ms
	return nil
}

func parseInitialOffset(value string) (initialOffset int64, err error) {
	switch {
	case strings.EqualFold(value, types.OffsetOldest):
		initialOffset = sarama.OffsetOldest
	case strings.EqualFold(value, types.OffsetNewest):
		initialOffset = sarama.OffsetNewest
	case value != "":
		return 0, fmt.Errorf("kafka error: invalid initialOffset: %s", value)
	default:
		initialOffset = sarama.OffsetNewest // Default
	}

	return initialOffset, err
}
