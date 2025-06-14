package sarama

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logx"
	ztrace "github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/core/utils"
	"github.com/zeromicro/go-zero/kafka/internal/metrics"
	"github.com/zeromicro/go-zero/kafka/internal/trace"
	"github.com/zeromicro/go-zero/kafka/internal/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	otrace "go.opentelemetry.io/otel/trace"
)

var (
	_             types.Producer = (*Producer)(nil)
	errEmptyTopic                = errors.New("topic is required for producer message")
)

type Producer struct {
	name         string
	appName      string
	defaultTopic string

	producer sarama.SyncProducer
	once     sync.Once
}

func NewProducer(pc types.ProducerConfig) (*Producer, error) { //nolint
	sc, err := toSaramaConfig(pc.Client)
	if err != nil {
		return nil, err
	}

	c := &Producer{
		name:         pc.Client.GetClientName(),
		defaultTopic: pc.Topic,
		appName:      pc.AppName,
	}

	producer, err := getSyncProducer(pc, *sc)
	if err != nil {
		return nil, err
	}

	c.producer = producer

	return c, nil
}

func newProducerFromClient(c *Client, defaultTopic string) (*Producer, error) {
	cc := &Producer{
		name:         c.config.Client.GetClientName(),
		defaultTopic: defaultTopic,
		appName:      c.config.Producer.AppName,
	}

	producer, err := sarama.NewSyncProducerFromClient(c.saramaClient)
	if err != nil {
		return nil, err
	}

	cc.producer = producer

	return cc, nil
}

// Send is low level API for generated code.
func (p *Producer) Send(ctx context.Context, messages ...*types.Message) error {
	tracer := otel.Tracer(ztrace.TracerName)
	msgs := make([]*sarama.ProducerMessage, 0, len(messages))
	spans := make([]otrace.Span, 0, len(messages))
	start := time.Now()

	for _, message := range messages {
		// compatible with old version
		if len(message.Topic) == 0 {
			message.Topic = p.defaultTopic
		}
		if len(message.Topic) == 0 {
			return errEmptyTopic
		}
		if len(message.Key) == 0 {
			message.Key = []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
		}

		attrs := injectHeaders(p.appName, message)
		attrs = append(attrs, attribute.KeyValue{
			Key:   types.KafkaBroker,
			Value: attribute.StringValue(p.name),
		}, attribute.KeyValue{
			Key:   types.KafkaTopic,
			Value: attribute.StringValue(message.Topic),
		})
		newCtx, span := tracer.Start(ctx, fmt.Sprintf("%s publish", message.Topic))
		span.SetAttributes(attrs...)
		spans = append(spans, span)

		// inject trace
		otel.GetTextMapPropagator().Inject(newCtx, trace.NewMessageCarrier(message))

		msgs = append(msgs, message2ProducerMessage(message))
	}

	err := p.producer.SendMessages(msgs)
	cost := time.Since(start).Seconds()

	for i, message := range messages {
		span := spans[i]
		if err != nil {
			span.RecordError(err)
			metrics.KafkaPublishCounter.Inc(metrics.ProducerTypeSync, p.name, message.Topic, metrics.CodeError)
		} else {
			// add partition and offset to input messages
			message.Partition = int(msgs[i].Partition)
			message.Offset = msgs[i].Offset
			metrics.KafkaPublishCounter.Inc(metrics.ProducerTypeSync, p.name, message.Topic, metrics.CodeOK)
			span.SetAttributes(
				semconv.MessagingMessageID(strconv.FormatInt(msgs[i].Offset, 10)),
				semconv.MessagingKafkaDestinationPartition(int(msgs[i].Partition)),
			)
		}
		metrics.KafkaPublishHistogram.ObserveFloat(cost,
			metrics.ProducerTypeSync, p.name, message.Topic)
		metrics.KafkaPayloadSizeHistogram.Observe(int64(len(message.Value)),
			message.Topic, fmt.Sprint(message.Partition))
		metrics.KafkaPublishGauge.Set(cost, metrics.ProducerKey{
			ProducerType: metrics.ProducerTypeSync,
			Brokers:      p.name,
			Topic:        message.Topic,
		})
		span.End()
	}
	return err
}

func (p *Producer) SendDelay(ctx context.Context, delaySeconds int64, messages ...*types.Message) error {
	if delaySeconds <= 5 {
		return errors.New("delaySeconds must>5")
	}
	if delaySeconds >= types.DelayTimeout {
		return errors.New("delaySeconds must be within 90 days")
	}
	if messages == nil || len(messages) == 0 {
		return errors.New("messages must not empty")
	}
	if len(messages) > 1 {
		for _, msg := range messages {
			if len(msg.Headers) > 0 {
				var newHeaders = make([]types.Header, 0)
				for _, h := range msg.Headers {
					newHeaders = append(newHeaders, h)
				}
				msg.Headers = newHeaders
			}
			msg.BuildDelayMessage(delaySeconds)
		}
	} else {
		messages[0].BuildDelayMessage(delaySeconds)
	}
	return p.Send(ctx, messages...)
}

func (p *Producer) Close() error {
	if p.producer == nil {
		return nil
	}

	var err error
	p.once.Do(func() {
		logx.Infof("kafka producer close, %s", p.name)
		// sarama producer close will close Client when we use sarama.NewSyncProducer,
		// will not close Client when use sarama.NewSyncProducerFromClient.
		err = p.producer.Close()
	})

	return err
}

func injectHeaders(appName string, message *types.Message) []attribute.KeyValue {
	nowTsStr := strconv.FormatInt(utils.CurrentMillis(), 10)
	clientID := appName + ".producer"

	message.SetHeader(types.ContentTypeKey, types.ContentTypeJSON)
	message.SetHeader(types.OriginAppNameKey, appName)
	message.SetHeader(types.ClientIDKey, clientID)
	message.SetHeader(types.CreateTimestampKey, nowTsStr)

	return []attribute.KeyValue{
		attribute.String(types.OriginAppNameKey, appName),
		attribute.String(types.ClientIDKey, clientID),
		attribute.String(types.CreateTimestampKey, nowTsStr),
	}
}

func getSyncProducer(pc types.ProducerConfig, sc sarama.Config) (sarama.SyncProducer, error) { //nolint
	if err := fillProducerConfig(&sc, pc.SharedProducerConfig); err != nil {
		return nil, err
	}

	// wrap tracing
	return sarama.NewSyncProducer(pc.Client.Brokers, &sc)
}

func fillProducerConfig(sc *sarama.Config, pc types.SharedProducerConfig) error {
	requiredAcks, err := requiredAcksFromString(pc.RequiredAcks)
	if err != nil {
		return err
	}

	// Add SyncProducer specific properties to copy of base config
	sc.Producer.RequiredAcks = requiredAcks
	sc.Producer.Retry.Max = 5
	sc.Producer.Return.Successes = true

	maxMessageBytes := pc.MaxMessageBytes
	if maxMessageBytes > 0 {
		sc.Producer.MaxMessageBytes = maxMessageBytes
	}

	// set compression from user config
	compression, err := compressionFromString(pc.Compression)
	if err != nil {
		return err
	}
	sc.Producer.Compression = compression

	// set partitioner from user config
	partitioner := partitionerFromString(pc.Partitioner)
	sc.Producer.Partitioner = partitioner

	if pc.Idempotent {
		sc.Producer.RequiredAcks = sarama.WaitForAll
		sc.Net.MaxOpenRequests = 1
		sc.Producer.Retry.Max = 5
	}

	sc.Producer.Flush.Bytes = pc.Flush.Bytes
	sc.Producer.Flush.Messages = pc.Flush.Messages
	sc.Producer.Flush.Frequency = pc.Flush.Frequency
	sc.Producer.Flush.MaxMessages = pc.Flush.MaxMessages
	return nil
}
