package sarama

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace/tracetest"
	"github.com/zeromicro/go-zero/internal/mock/saramamock"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

//go:generate mockgen -destination=../../../internal/mock/saramamock/sarama.gen.go -package=saramamock github.com/IBM/sarama PartitionConsumer,Consumer,Client,AsyncProducer

func Test_partitionConsumer_handleMessage(t *testing.T) {
	t.Run("should recover user panic", func(t *testing.T) {
		logx.Disable()
		p := &partitionConsumer{
			consumer: &consumer{
				name: "test",
			},
			handler: func(ctx context.Context, msg *types.Message) {
				panic("user panic")
			},
		}
		assert.NotPanics(t, func() {
			p.handleMessage(&sarama.ConsumerMessage{})
		})
	})

	t.Run("should add metrics and trace", func(t *testing.T) {
		me := tracetest.NewInMemoryExporter(t)
		done := make(chan struct{}, 1)
		p := &partitionConsumer{
			consumer: &consumer{
				name: "test",
			},
			handler: func(ctx context.Context, msg *types.Message) {
				done <- struct{}{}
			},
		}

		p.handleMessage(&sarama.ConsumerMessage{
			Key:       []byte("test"),
			Value:     []byte("test data"),
			Topic:     "abc",
			Partition: 1,
			Offset:    10,
		})
		<-done
		assert.True(t, len(me.GetSpans()) > 0)
		assert.Equal(t, "kafka.consumer.handler", me.GetSpans()[0].Name)
	})
}

func Test_consumer_ConsumePartition(t *testing.T) {
	logx.Disable()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sc := saramamock.NewMockConsumer(ctrl)

		spc := saramamock.NewMockPartitionConsumer(ctrl)
		messages := make(chan *sarama.ConsumerMessage, 10)
		errs := make(chan *sarama.ConsumerError, 10)
		spc.EXPECT().Messages().MinTimes(1).Return(messages)
		spc.EXPECT().Errors().MinTimes(1).Return(errs)
		spc.EXPECT().AsyncClose().MinTimes(1)

		sc.EXPECT().ConsumePartition("test", int32(0), int64(1)).
			MinTimes(1).Return(spc, nil)

		c := &consumer{
			name:           "test111",
			saramaConsumer: sc,
		}

		done := make(chan struct{})
		i := 0
		errCount := int32(0)
		cp, err := c.ConsumePartition("test", 0, 1, func(ctx context.Context, message *types.Message) {
			i++
			if i == 3 {
				done <- struct{}{}
			}
		}, func(consumerError *ConsumerError) {
			atomic.AddInt32(&errCount, 1)
		})
		assert.NoError(t, err)

		messages <- &sarama.ConsumerMessage{
			Topic: "test",
			Value: []byte("1"),
		}

		errs <- &sarama.ConsumerError{
			Topic:     "test",
			Partition: 0,
			Err:       errors.New("test"),
		}

		messages <- &sarama.ConsumerMessage{
			Topic: "test",
			Value: []byte("2"),
		}

		messages <- &sarama.ConsumerMessage{
			Topic: "test",
			Value: []byte("3"),
		}

		<-done
		assert.NoError(t, cp.Close())
		assert.Equal(t, 3, i)
		assert.Equal(t, int32(1), atomic.LoadInt32(&errCount))
	})

	t.Run("no error handler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sc := saramamock.NewMockConsumer(ctrl)

		spc := saramamock.NewMockPartitionConsumer(ctrl)
		messages := make(chan *sarama.ConsumerMessage, 10)
		errs := make(chan *sarama.ConsumerError, 10)
		spc.EXPECT().Messages().MinTimes(1).Return(messages)
		spc.EXPECT().Errors().MinTimes(1).Return(errs)
		spc.EXPECT().AsyncClose().MinTimes(1)

		sc.EXPECT().ConsumePartition("test", int32(0), int64(1)).
			MinTimes(1).Return(spc, nil)

		c := &consumer{
			name:           "test111",
			saramaConsumer: sc,
		}

		done := make(chan struct{})
		i := 0
		cp, err := c.ConsumePartition("test", 0, 1, func(ctx context.Context, message *types.Message) {
			i++
			if i == 3 {
				done <- struct{}{}
			}
		}, nil)
		assert.NoError(t, err)

		messages <- &sarama.ConsumerMessage{
			Topic: "test",
			Value: []byte("1"),
		}

		errs <- &sarama.ConsumerError{
			Topic:     "test",
			Partition: 0,
			Err:       errors.New("test"),
		}

		messages <- &sarama.ConsumerMessage{
			Topic: "test",
			Value: []byte("2"),
		}

		messages <- &sarama.ConsumerMessage{
			Topic: "test",
			Value: []byte("3"),
		}

		<-done
		assert.NoError(t, cp.Close())
		assert.Equal(t, 3, i)
	})

	t.Run("underline ConsumePartition error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sc := saramamock.NewMockConsumer(ctrl)

		sc.EXPECT().ConsumePartition("test", int32(0), int64(1)).
			MinTimes(1).Return(nil, errors.New("test error"))

		c := &consumer{
			name:           "test111",
			saramaConsumer: sc,
		}

		_, err := c.ConsumePartition("test", 0, 1, nil, nil)
		assert.EqualError(t, err, "test error")
	})
}

func Test_partitionConsumer_HighWaterMarkOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	spc := saramamock.NewMockPartitionConsumer(ctrl)
	spc.EXPECT().HighWaterMarkOffset().Return(int64(10))
	p := &partitionConsumer{
		spc: spc,
	}
	assert.Equal(t, int64(10), p.HighWaterMarkOffset())
}

func Test_consumer_Close(t *testing.T) {
	logx.Disable()

	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sc := saramamock.NewMockClient(ctrl)
		sc.EXPECT().Close().MinTimes(1).MaxTimes(1).Return(nil)
		c := &consumer{
			saramaClient: sc,
			closeClient:  true,
		}
		assert.NoError(t, c.Close())
		assert.NoError(t, c.Close())
	})

	t.Run("should close children partitionConsumers", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sc := saramamock.NewMockClient(ctrl)
		sc.EXPECT().Close().MinTimes(1).MaxTimes(1).Return(nil)
		spc1 := saramamock.NewMockPartitionConsumer(ctrl)
		spc1.EXPECT().AsyncClose().MinTimes(1)
		spc2 := saramamock.NewMockPartitionConsumer(ctrl)
		spc2.EXPECT().AsyncClose().MinTimes(1)

		c := &consumer{
			saramaClient: sc,
			partitionConsumers: []*partitionConsumer{
				&partitionConsumer{
					spc:  spc1,
					done: make(chan struct{}),
				},
				&partitionConsumer{
					spc:  spc2,
					done: make(chan struct{}),
				},
			},
		}
		assert.NoError(t, c.Close())
		assert.NoError(t, c.Close())
	})
}

func Test_consumer_GetNewestOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	sc := saramamock.NewMockClient(ctrl)
	sc.EXPECT().GetOffset("test", int32(0), sarama.OffsetNewest).Return(int64(10), nil)
	c := &consumer{
		saramaClient: sc,
	}
	offset, err := c.GetNewestOffset("test", 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), offset)
}

func TestNewConsumer(t *testing.T) {
	t.Run("invalid params", func(t *testing.T) {
		_, err := NewConsumer(types.ConsumerConfig{})
		assert.Error(t, err)
	})

	t.Run("createClient fauk", func(t *testing.T) {
		_, err := NewConsumer(types.ConsumerConfig{
			Client: types.ClientConfig{Brokers: []string{"b1"}},
		})
		assert.Error(t, err)
	})
}
