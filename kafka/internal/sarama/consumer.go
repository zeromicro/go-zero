package sarama

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/kafka/internal/metrics"
	ztrace "github.com/zeromicro/go-zero/kafka/internal/trace"
	"github.com/zeromicro/go-zero/kafka/internal/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otrace "go.opentelemetry.io/otel/trace"
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
	PartitionConsumer interface {
		// Close closes the partition consumer,
		// the underline sarama.PartitionConsumer.AsyncClose will be called,
		// and exit the goroutine that is consuming messages and errors.
		Close() error
		// HighWaterMarkOffset returns the high water mark offset of the partition,
		// i.e. the offset that will be used for the next message that will be produced.
		// You can use this to determine how far behind the processing is.
		HighWaterMarkOffset() int64
	}

	// ConsumerHandler is a handler for consuming messages.
	ConsumerHandler func(ctx context.Context, message *types.Message)
	// ConsumerError is type alias of sarama.ConsumerError.
	ConsumerError = sarama.ConsumerError
	// ConsumerErrorHandler is a handler for consuming errors.
	ConsumerErrorHandler func(consumerError *ConsumerError)

	consumer struct {
		name               string
		consumerConfig     types.ConsumerConfig
		saramaConsumer     sarama.Consumer
		saramaClient       sarama.Client
		partitionConsumers []*partitionConsumer
		once               sync.Once
		// 老版本 consumer 会自己创建 saramaClient, 所以需要手动关闭
		closeClient bool

		sharedClient *Client
	}

	partitionConsumer struct {
		consumer     *consumer
		spc          sarama.PartitionConsumer
		handler      ConsumerHandler
		errorHandler ConsumerErrorHandler
		done         chan struct{}
		once         sync.Once
	}
)

func NewConsumer(config types.ConsumerConfig) (Consumer, error) {
	c := &consumer{
		name:        config.Client.GetClientName(),
		closeClient: true,
	}

	sc, err := toSaramaConfig(config.Client)
	if err != nil {
		return nil, err
	}

	if err := fillConsumerConfig(sc, config.SharedConsumerConfig); err != nil {
		return nil, err
	}

	sClient, err := sarama.NewClient(config.Client.Brokers, sc)
	if err != nil {
		return nil, err
	}
	sConsumer, err := sarama.NewConsumerFromClient(sClient)
	if err != nil {
		return nil, err
	}
	c.saramaClient = sClient
	c.saramaConsumer = sConsumer

	return c, nil
}

func newConsumerFromClient(client *Client) (*consumer, error) {
	c := &consumer{
		name: client.config.Client.GetClientName(),
	}
	sConsumer, err := sarama.NewConsumerFromClient(client.saramaClient)
	if err != nil {
		return nil, err
	}

	c.saramaClient = client.saramaClient
	c.saramaConsumer = sConsumer

	return c, nil
}

// ConsumePartition warps sarama.consumer.ConsumePartition with handlers,
// errorHandler is optional, if it is nil, the error will be logged.
func (c *consumer) ConsumePartition(
	topic string, partition int32, offset int64,
	handler ConsumerHandler, errorHandler ConsumerErrorHandler,
) (PartitionConsumer, error) {
	pc, err := c.saramaConsumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		return nil, err
	}

	p := &partitionConsumer{
		consumer:     c,
		spc:          pc,
		handler:      handler,
		errorHandler: errorHandler,
		done:         make(chan struct{}),
	}

	p.setupCallbacks()
	c.partitionConsumers = append(c.partitionConsumers, p)

	logx.Infof("consumer partition started, name: %s, topic: %s, partition: %d, startOffset: %d",
		p.consumer.name, topic, partition, offset)

	return p, nil
}

// GetNewestOffset returns the newest offset of the partition.
func (c *consumer) GetNewestOffset(topic string, partition int32) (int64, error) {
	return c.saramaClient.GetOffset(topic, partition, sarama.OffsetNewest)
}

func (c *consumer) Close() error {
	if c.saramaClient == nil {
		return nil
	}

	var err error
	c.once.Do(func() {
		logx.Info("consumer close")
		for _, p := range c.partitionConsumers {
			// close the partition consumer
			_ = p.Close()
		}
		if c.closeClient {
			err = c.saramaClient.Close()
		}
	})

	return err
}

func (p *partitionConsumer) setupCallbacks() {
	go func() {
		for {
			select {
			case <-p.done:
				return
			case msg, ok := <-p.spc.Messages():
				if !ok {
					return
				}
				p.handleMessage(msg)
			case err := <-p.spc.Errors():
				if err != nil {
					if p.errorHandler != nil {
						threading.RunSafe(func() {
							p.errorHandler(err)
						})
					} else {
						logx.Errorw(fmt.Sprintf("consumer error: %v",
							err), logx.Field("brokers", p.consumer.name))
					}
				}
			}
		}
	}()
}

func (p *partitionConsumer) handleMessage(message *sarama.ConsumerMessage) {
	innerHandleMessage(p.consumer, message, false, func(ctx context.Context, msg *types.Message) bool {
		p.handler(ctx, msg)
		return true
	})
}

func innerHandleMessage(c *consumer, message *sarama.ConsumerMessage, isUnsafe bool, f func(ctx context.Context, message *types.Message) bool) bool {
	continueConsuming := true
	start := time.Now()
	// extract span info from consumer message into ctx.
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(message))

	trace.DoWithSpan(ctx, "kafka.consumer.handler", func(ctx context.Context, span otrace.Span) {
		msg := consumerMessage2Message(message)
		ztrace.AddConsumerAttrs(span, msg, attribute.KeyValue{
			Key:   types.KafkaBroker,
			Value: attribute.StringValue(c.name),
		})
		if isUnsafe {
			continueConsuming = f(ctx, msg)
		} else {
			threading.RunSafe(func() {
				continueConsuming = f(ctx, msg)
			})
		}
		cost := time.Since(start)
		partition := strconv.Itoa(int(message.Partition))
		metrics.KafkaConsumerCounter.Inc(c.name, message.Topic, partition)
		metrics.KafkaConsumerHistogram.ObserveFloat(cost.Seconds(), c.name,
			message.Topic, partition)
		metrics.KafkaConsumerGauge.Set(cost.Seconds(), metrics.ConsumerKey{
			Brokers:   c.name,
			Topic:     message.Topic,
			Partition: partition,
		})

	}, otrace.WithSpanKind(otrace.SpanKindConsumer))

	return continueConsuming
}

func (p *partitionConsumer) Close() error {
	if p.spc == nil {
		return nil
	}

	p.once.Do(func() {
		logx.Info("partitionConsumer close")
		close(p.done)
		p.spc.AsyncClose()
	})

	return nil
}

func (p *partitionConsumer) HighWaterMarkOffset() int64 {
	return p.spc.HighWaterMarkOffset()
}
