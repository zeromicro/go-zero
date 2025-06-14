package sarama

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/internal/mock/saramamock"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

//go:generate mockgen -destination=../../../internal/mock/saramamock/sarama.gen.go -package=saramamock github.com/Shopify/sarama PartitionConsumer,Consumer,Client,AsyncProducer

type msgReceiver struct {
	i               *int
	errCount        *int
	LoopInitCount   *int
	LoopFinishCount *int
	d               chan struct{}
}

func (c *msgReceiver) OnKNameConfirmed(k string) {
}

// 执行完实际业务处理后再调用这个方法判断
func (c *msgReceiver) OnConsumeLoopInit(currentBeginInclusive, offset, currentEndExclusive int64) bool {
	*c.LoopInitCount++
	return true
}

// 执行完实际业务处理后再调用这个方法判断
func (c *msgReceiver) OnConsumeMessage(ctx context.Context, msg *types.Message, highWaterMarkOffset int64) bool {
	*c.i++
	if *c.i == 3 {
		c.d <- struct{}{}
	}
	return true
}

func (c *msgReceiver) OnConsumeError(consumerErr *ConsumerError, highWaterMarkOffset int64) (continueConsuming bool) {
	*c.errCount++
	return true
}

func (c *msgReceiver) OnConsumeLoopFinish(finishReason FinishReason) {
	*c.LoopFinishCount++
}

func Test_unsafeConsumer_SpawnPartitionConsumer(t *testing.T) {
	logx.Disable()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sc := saramamock.NewMockConsumer(ctrl)

		spc := saramamock.NewMockPartitionConsumer(ctrl)
		mc := saramamock.NewMockClient(ctrl)

		messages := make(chan *sarama.ConsumerMessage, 10)
		errs := make(chan *sarama.ConsumerError, 10)
		spc.EXPECT().Messages().MinTimes(1).Return(messages)
		spc.EXPECT().Errors().MinTimes(1).Return(errs)
		spc.EXPECT().Close().MinTimes(1)
		spc.EXPECT().HighWaterMarkOffset().MinTimes(1).Return(int64(10))

		sc.EXPECT().ConsumePartition("test", int32(0), int64(1)).MinTimes(1).Return(spc, nil)

		sharedClient := &Client{
			saramaClient: mc,
		}

		c := &consumer{
			name:           "test111",
			saramaConsumer: sc,
			sharedClient:   sharedClient,
		}
		uc := &unsafeConsumer{
			consumer: c,
			k:        "k123",
		}

		done := make(chan struct{})
		i := 0
		errCount := 0
		LoopInitCount := 0
		LoopFinishCount := 0

		handler0 := &msgReceiver{&i, &errCount, &LoopInitCount, &LoopFinishCount, done}
		cp, err := uc.SpawnPartitionConsumer("test", 0, 1, 0, 2, handler0)
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
		assert.Equal(t, 1, errCount)
		assert.Equal(t, 1, LoopInitCount)
	})
}
func Test_unsafeConsumer_SameUniq(t *testing.T) {
	logx.Disable()

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sc := saramamock.NewMockConsumer(ctrl)

		spc := saramamock.NewMockPartitionConsumer(ctrl)
		mc := saramamock.NewMockClient(ctrl)

		messages := make(chan *sarama.ConsumerMessage, 10)
		errs := make(chan *sarama.ConsumerError, 10)
		spc.EXPECT().Messages().MaxTimes(10).Return(messages)
		spc.EXPECT().Errors().MaxTimes(1).Return(errs)
		sc.EXPECT().ConsumePartition("test", int32(0), int64(1)).MinTimes(1).Return(spc, nil)

		sharedClient := &Client{
			saramaClient: mc,
		}

		c := &consumer{
			name:           "test111",
			saramaConsumer: sc,
			sharedClient:   sharedClient,
		}
		uc := &unsafeConsumer{
			consumer: c,
			k:        "k123",
		}

		done := make(chan struct{})
		i := 0
		errCount := 0
		LoopInitCount := 0
		LoopFinishCount := 0

		handler0 := &msgReceiver{&i, &errCount, &LoopInitCount, &LoopFinishCount, done}
		_, err := uc.SpawnPartitionConsumer("test", 0, 1, 0, 2, handler0)
		assert.NoError(t, err)

		_, err = uc.SpawnPartitionConsumer("test", 0, 1, 0, 2, handler0)
		assert.Error(t, err)
	})
}

func Test_unsafeConsumer_Close(t *testing.T) {
	logx.Disable()

	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sc := saramamock.NewMockClient(ctrl)
		sc.EXPECT().Close().MinTimes(1).MaxTimes(1).Return(nil)
		spc1 := saramamock.NewMockPartitionConsumer(ctrl)
		spc1.EXPECT().AsyncClose().MinTimes(1)

		sharedClient := &Client{
			saramaClient:              sc,
			spawnedPartitionConsumers: sync.Map{},
		}

		c := &consumer{
			saramaClient: sc,
			closeClient:  true,
			sharedClient: sharedClient,
		}

		uc := &unsafeConsumer{
			consumer: c,
			k:        "k123",
		}

		usp := &unsafePartitionConsumer{
			unsafeConsumer: uc,
			spc:            spc1,
			done:           make(chan struct{}),
		}
		sharedClient.spawnedPartitionConsumers.Store("", usp)
		assert.NoError(t, uc.Close())
		assert.NoError(t, uc.Close())
	})
}

func Test_unsafeConsumer_ClosePartitionConsumer(t *testing.T) {
	logx.Disable()

	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sc := saramamock.NewMockClient(ctrl)
		sc.EXPECT().Close().MinTimes(1).MaxTimes(1).Return(nil)
		spc1 := saramamock.NewMockPartitionConsumer(ctrl)
		spc1.EXPECT().AsyncClose().MinTimes(1)

		sharedClient := &Client{
			saramaClient:              sc,
			spawnedPartitionConsumers: sync.Map{},
		}

		c := &consumer{
			saramaClient: sc,
			closeClient:  true,
			sharedClient: sharedClient,
		}

		uc := &unsafeConsumer{
			consumer: c,
			k:        "k123",
		}

		usp := &unsafePartitionConsumer{
			unsafeConsumer: uc,
			spc:            spc1,
			done:           make(chan struct{}),
		}
		sharedClient.spawnedPartitionConsumers.Store("test-0", usp)
		assert.Equal(t, true, uc.ClosePartitionConsumer(context.Background(), "test", 0))
	})
}

func Test_unsafeConsumer_SyncClosePartitionConsumer(t *testing.T) {
	logx.Disable()

	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sc := saramamock.NewMockClient(ctrl)
		sc.EXPECT().Close().MinTimes(1).MaxTimes(1).Return(nil)
		spc1 := saramamock.NewMockPartitionConsumer(ctrl)
		spc1.EXPECT().Close().MinTimes(1)

		sharedClient := &Client{
			saramaClient:              sc,
			spawnedPartitionConsumers: sync.Map{},
		}

		c := &consumer{
			saramaClient: sc,
			closeClient:  true,
			sharedClient: sharedClient,
		}

		uc := &unsafeConsumer{
			consumer: c,
			k:        "k123",
		}

		usp := &unsafePartitionConsumer{
			unsafeConsumer: uc,
			spc:            spc1,
			done:           make(chan struct{}),
		}
		sharedClient.spawnedPartitionConsumers.Store("test-0", usp)
		assert.Equal(t, true, uc.SyncClosePartitionConsumer(context.Background(), "test", 0))
	})
}
