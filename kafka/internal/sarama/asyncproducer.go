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
	"github.com/zeromicro/go-zero/core/threading"
	ztrace "github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/kafka/internal/metrics"
	"github.com/zeromicro/go-zero/kafka/internal/trace"
	"github.com/zeromicro/go-zero/kafka/internal/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	otrace "go.opentelemetry.io/otel/trace"
)

type (
	// AsyncProducer wrapper sarama async producer.
	AsyncProducer struct {
		name         string
		appName      string
		defaultTopic string
		producer     sarama.AsyncProducer

		opt     AsyncProducerOption
		closeCh chan struct{}
		once    sync.Once

		// NewAsyncProducer API 会自己创建 client，
		// Client API 创建的则共享 client。
		sharedClient *Client
	}

	// AsyncProducerOption is option for AsyncProducer.
	AsyncProducerOption struct {
		SuccessHandler SuccessHandler
		ErrorHandler   ErrorHandler
	}

	// ProducerError is like sarama.ProducerError,
	// but Msg type is *types.Message.
	ProducerError struct {
		Err error
		Msg *types.Message
	}

	// SuccessHandler is optional callback function for handling success message
	// from sarama async produce Successes() channel.
	SuccessHandler func(msg *types.Message)
	// ErrorHandler is optional callback function for handling error message
	// from sarama async produce Errors() channel.
	ErrorHandler func(producerError *ProducerError)

	// CallbackHandler is per message scope callback handler.
	CallbackHandler func(msg *types.Message, err error)

	producerMetadata struct {
		startTime       time.Time
		callbackHandler CallbackHandler
		span            otrace.Span
	}
)

func (pe ProducerError) Error() string {
	return fmt.Sprintf("kafka: Failed to produce message to topic %s: %s", pe.Msg.Topic, pe.Err)
}

func (pe ProducerError) Unwrap() error {
	return pe.Err
}

type ProducerErrors []*ProducerError

func (pe ProducerErrors) Error() string {
	return fmt.Sprintf("kafka: Failed to deliver %d messages.", len(pe))
}

// NewAsyncProducer create an AsyncProducer instance.
func NewAsyncProducer(pc types.ProducerConfig, opt AsyncProducerOption) (*AsyncProducer, error) { //nolint
	sc, err := toSaramaConfig(pc.Client)
	if err != nil {
		return nil, err
	}

	c := &AsyncProducer{
		name:         pc.Client.GetClientName(),
		defaultTopic: pc.Topic,
		appName:      pc.AppName,
		closeCh:      make(chan struct{}),
		opt:          opt,
		sharedClient: &Client{isExiting: true}, // 主要防止 sharedClient.pwg 空指针, 并且不需要退出判断 panic
	}

	producer, err := createAsyncProducer(pc, *sc)
	if err != nil {
		return nil, err
	}

	c.producer = producer
	c.setupCallbacks()

	return c, nil
}

func newAsyncProducerFromClient(c *Client, defaultTopic string, opt AsyncProducerOption) (*AsyncProducer, error) {
	cc := &AsyncProducer{
		name:         c.config.Client.GetClientName(),
		defaultTopic: defaultTopic,
		appName:      c.config.Producer.AppName,
		closeCh:      make(chan struct{}),
		opt:          opt,
		sharedClient: c,
	}

	producer, err := sarama.NewAsyncProducerFromClient(c.saramaClient)
	if err != nil {
		return nil, err
	}

	// wrap tracing
	// todo: 暂时不要使用 otelsarama.WrapAsyncProducer, 他有个全局缓存 map<spanID, metadata>,
	// 如果我们没开启 trace, spanID 为空时, 我们的 metadata 会被修改为一个空 SpanID.
	//cc.producer = otelsarama.WrapAsyncProducer(c.sc, producer)
	cc.producer = producer
	return cc, nil
}

// Send sends messages to kafka async producer,
// this API signature is same as sync producer,
// but it always returns nil error immediately,
// you should use WithSuccessHandler and WithErrorHandler to set success/error handler.
func (a *AsyncProducer) Send(ctx context.Context, messages ...*types.Message) error {
	return a.SendWithCallback(ctx, messages, nil)
}

// SendWithCallback sends messages to kafka async producer with per message scope callback handler.
func (a *AsyncProducer) SendWithCallback(ctx context.Context, messages []*types.Message, handler CallbackHandler) error {
	return a.sendWithCallback(ctx, messages, func(sMsg *sarama.ProducerMessage, message *types.Message, span otrace.Span) {
		// sarama message metadata 可以保存传递元数据, 这些信息会在 callback 消息中拿到,
		// 并不会发送出去, 所以我们使用这个字段来保存发送时间, 用于计算延迟,
		// 这样可以避免 header 中的时间戳被篡改的问题, 也不需要类型转换.
		setInternalMetadata(sMsg, &producerMetadata{
			startTime:       time.Now(),
			callbackHandler: handler,
			span:            span,
		})
	})
}

func (a *AsyncProducer) sendWithCallback(ctx context.Context, messages []*types.Message, innerFunc func(sMsg *sarama.ProducerMessage, message *types.Message, span otrace.Span)) error {
	tracer := otel.Tracer(ztrace.TracerName)
	for _, message := range messages {
		// compatible with old version
		if len(message.Topic) == 0 {
			message.Topic = a.defaultTopic
		}
		if len(message.Topic) == 0 {
			return errEmptyTopic
		}
		if len(message.Key) == 0 {
			message.Key = []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
		}

		attrs := injectHeaders(a.appName, message)
		attrs = append(attrs, attribute.KeyValue{
			Key:   types.KafkaBroker,
			Value: attribute.StringValue(a.name),
		}, attribute.KeyValue{
			Key:   types.KafkaTopic,
			Value: attribute.StringValue(message.Topic),
		})

		newCtx, span := tracer.Start(ctx, fmt.Sprintf("%s publish", message.Topic))
		span.SetAttributes(attrs...)

		// inject trace
		otel.GetTextMapPropagator().Inject(newCtx, trace.NewMessageCarrier(message))
		sMsg := message2ProducerMessage(message)
		innerFunc(sMsg, message, span)
		a.producer.Input() <- sMsg
	}

	return nil
}

func (a *AsyncProducer) SendDelay(ctx context.Context, delaySeconds int64, messages ...*types.Message) error {
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
	return a.Send(ctx, messages...)
}

// Close closes underlying async producer and close background callbacks loop.
func (a *AsyncProducer) Close() error {
	if a.producer == nil {
		return nil
	}

	var err error
	a.once.Do(func() {
		close(a.closeCh)
		logx.Infof("kafka async producer close, %s", a.name)
		// sarama producer close will close Client when we use sarama.NewAsyncProducer,
		// will not close Client when use sarama.NewAsyncProducerFromClient.
		err = a.producer.Close()
	})

	return err
}

// setupCallbacks setup callbacks loop for async producer.
func (a *AsyncProducer) setupCallbacks() {
	a.setupCallbacksRunSafe(a.sharedClient.config.Producer.EnableRecovery)
}

func (a *AsyncProducer) setupCallbacksRunSafe(runSafe bool) {
	a.sharedClient.pwg.Add(1) // 新创建一个producer的消费协程, 引用计数+1
	go func() {
		defer a.sharedClient.pwg.Done() // ProcessorLoop协程结束后, 引用计数-1
	SharedProducerResultCallbackProcessorLoop:
		for {
			select {
			case <-a.closeCh:
				return
			case msg, ok := <-a.producer.Successes():
				if !ok {
					logx.Infof("SharedProducer_%s Successes channel closed", a.name)
					// sarama底层是先close(errors) 再close(successes)
					break SharedProducerResultCallbackProcessorLoop
				}
				a.reportMetrics(msg, metrics.CodeOK)

				md := getInternalMetadata(msg).(*producerMetadata)
				if md.callbackHandler != nil {
					if !runSafe {
						md.callbackHandler(producerMessage2Message(msg), nil)
					} else {
						threading.RunSafe(func() {
							md.callbackHandler(producerMessage2Message(msg), nil)
						})
					}
				} else if a.opt.SuccessHandler != nil {
					if !runSafe {
						a.opt.SuccessHandler(producerMessage2Message(msg))
					} else {
						threading.RunSafe(func() {
							a.opt.SuccessHandler(producerMessage2Message(msg))
						})
					}
				}

				md.span.SetAttributes(
					semconv.MessagingMessageID(strconv.FormatInt(msg.Offset, 10)),
					semconv.MessagingKafkaDestinationPartition(int(msg.Partition)),
				)
				md.span.End()
			}
		}
		// producer loop结束后, 判断下是用户发起的close 还是底层error触发的close 决定是否直接panic
		if !a.sharedClient.isExiting {
			panic(fmt.Errorf("unexpected SharedProducer_%s producer loop end", a.name))
		}
	}()

	go func() {
		for {
			select {
			case <-a.closeCh:
				return
			case errMsg, ok := <-a.producer.Errors():
				if !ok {
					return
				}
				a.reportMetrics(errMsg.Msg, metrics.CodeError)
				md := getInternalMetadata(errMsg.Msg).(*producerMetadata)
				if md.callbackHandler != nil {
					if !runSafe {
						md.callbackHandler(producerMessage2Message(errMsg.Msg), errMsg.Err)
					} else {
						threading.RunSafe(func() {
							md.callbackHandler(producerMessage2Message(errMsg.Msg), errMsg.Err)
						})
					}
				} else if a.opt.ErrorHandler != nil {
					if !runSafe {
						a.opt.ErrorHandler(&ProducerError{
							Err: errMsg.Err,
							Msg: producerMessage2Message(errMsg.Msg),
						})
					} else {
						threading.RunSafe(func() {
							a.opt.ErrorHandler(&ProducerError{
								Err: errMsg.Err,
								Msg: producerMessage2Message(errMsg.Msg),
							})
						})
					}
				}

				md.span.RecordError(errMsg.Err)
				md.span.End()
			}
		}
	}()
}

func (a *AsyncProducer) reportMetrics(msg *sarama.ProducerMessage, code string) {
	metrics.KafkaPublishCounter.Inc(metrics.ProducerTypeASync, a.name, msg.Topic, code)
	if metadata, ok := getInternalMetadata(msg).(*producerMetadata); ok {
		cost := time.Since(metadata.startTime).Seconds()
		metrics.KafkaPublishHistogram.ObserveFloat(cost,
			metrics.ProducerTypeASync, a.name, msg.Topic)
		metrics.KafkaPayloadSizeHistogram.Observe(int64(msg.Value.Length()),
			msg.Topic, fmt.Sprint(msg.Partition))
		metrics.KafkaPublishGauge.Set(cost, metrics.ProducerKey{
			ProducerType: metrics.ProducerTypeASync,
			Brokers:      a.name,
			Topic:        msg.Topic,
		})
	}
}

func createAsyncProducer(pc types.ProducerConfig, sc sarama.Config) (sarama.AsyncProducer, error) { //nolint
	if err := fillProducerConfig(&sc, pc.SharedProducerConfig); err != nil {
		return nil, err
	}

	producer, err := sarama.NewAsyncProducer(pc.Client.Brokers, &sc)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
