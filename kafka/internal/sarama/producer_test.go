package sarama

import (
	"context"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/internal/mock/saramamock"
	"github.com/zeromicro/go-zero/kafka/internal/types"
)

func Test_injectHeaders(t *testing.T) {
	msg := &types.Message{
		Topic: "test",
		Key:   []byte("abc"),
		Value: []byte("value"),
		Headers: []types.Header{
			{
				Key:   "test",
				Value: []byte("test"),
			},
		},
	}

	attrs := injectHeaders("test-app", msg)

	assert.Equal(t, 5, len(msg.Headers))
	assert.Equal(t, 3, len(attrs))

	for _, h := range msg.Headers {
		switch h.Key {
		case "test":
			assert.Equal(t, "test", string(h.Value))
		case types.ContentTypeKey:
			assert.Equal(t, "application/json", string(h.Value))
		case types.OriginAppNameKey:
			assert.Equal(t, "test-app", string(h.Value))
		case types.ClientIDKey:
			assert.Equal(t, "test-app.producer", string(h.Value))
		case types.CreateTimestampKey:
			assert.NotEmpty(t, h.Value)
		}
	}
}

func Test_fillProducerConfig(t *testing.T) {
	sc := &sarama.Config{}
	pc := types.SharedProducerConfig{
		RequiredAcks: "all",
		Idempotent:   true,
		Flush: struct {
			Bytes       int           `json:",default=65536"`
			Messages    int           `json:",default=0"`
			Frequency   time.Duration `json:",default=1ms"`
			MaxMessages int           `json:",default=100000"`
		}(struct {
			Bytes       int
			Messages    int
			Frequency   time.Duration
			MaxMessages int
		}{Bytes: 100, Messages: 100, Frequency: 0, MaxMessages: 100}),
	}

	err := fillProducerConfig(sc, pc)
	assert.NoError(t, err)

	assert.Equal(t, sc.Producer.RequiredAcks, sarama.WaitForAll)
	assert.Equal(t, sc.Net.MaxOpenRequests, 1)
	assert.Equal(t, sc.Producer.Retry.Max, 5)

	assert.Equal(t, sc.Producer.Flush.Bytes, 100)
	assert.Equal(t, sc.Producer.Flush.Messages, 100)
	assert.Equal(t, sc.Producer.Flush.Frequency, time.Duration(0))
	assert.Equal(t, sc.Producer.Flush.MaxMessages, 100)
}

func Test_fillPartitioner(t *testing.T) {
	sc := &sarama.Config{}
	pc := types.SharedProducerConfig{
		RequiredAcks: "all",
		Idempotent:   true,
		Flush: struct {
			Bytes       int           `json:",default=65536"`
			Messages    int           `json:",default=0"`
			Frequency   time.Duration `json:",default=1ms"`
			MaxMessages int           `json:",default=100000"`
		}(struct {
			Bytes       int
			Messages    int
			Frequency   time.Duration
			MaxMessages int
		}{Bytes: 100, Messages: 100, Frequency: 0, MaxMessages: 100}),
	}

	var partitioner sarama.Partitioner

	pc.Partitioner = "hash"
	fillProducerConfig(sc, pc)
	partitioner = sc.Producer.Partitioner("")
	assert.IsType(t, partitioner, sarama.NewHashPartitioner(""))

	pc.Partitioner = "roundrobin"
	fillProducerConfig(sc, pc)
	partitioner = sc.Producer.Partitioner("")
	assert.IsType(t, partitioner, sarama.NewRoundRobinPartitioner(""))

	pc.Partitioner = "random"
	fillProducerConfig(sc, pc)
	partitioner = sc.Producer.Partitioner("")
	assert.IsType(t, partitioner, sarama.NewRandomPartitioner(""))

	pc.Partitioner = "manual"
	fillProducerConfig(sc, pc)
	partitioner = sc.Producer.Partitioner("")
	assert.IsType(t, partitioner, sarama.NewManualPartitioner(""))
}

func TestProducer_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	ap := saramamock.NewMockSyncProducer(ctrl)
	p := &Producer{
		producer:     ap,
		defaultTopic: "xxx",
	}

	t.Run("send fail", func(t *testing.T) {
		ap.EXPECT().SendMessages(gomock.Any()).Return(assert.AnError).Times(1)
		err := p.Send(context.Background(), &types.Message{})
		assert.Error(t, err)
	})

	t.Run("send success", func(t *testing.T) {
		ap.EXPECT().SendMessages(gomock.Any()).Return(nil).Times(1)
		err := p.Send(context.Background(), &types.Message{})
		assert.NoError(t, err)
	})
}

func TestProducer_SendDelay(t *testing.T) {
	ctrl := gomock.NewController(t)
	ap := saramamock.NewMockSyncProducer(ctrl)
	p := &Producer{
		producer:     ap,
		defaultTopic: "xxx",
	}

	t.Run("send fail", func(t *testing.T) {
		ap.EXPECT().SendMessages(gomock.Any()).Return(assert.AnError).Times(1)
		err := p.SendDelay(context.Background(), 10, &types.Message{})
		assert.Error(t, err)
		err = p.SendDelay(context.Background(), 1, &types.Message{})
		assert.Error(t, err)
		err = p.SendDelay(context.Background(), 1+types.DelayTimeout, &types.Message{})
		assert.Error(t, err)

		err = p.SendDelay(context.Background(), 10)
		assert.Error(t, err)
	})

	t.Run("send success", func(t *testing.T) {
		ap.EXPECT().SendMessages(gomock.Any()).Return(nil).Times(2)
		err := p.SendDelay(context.Background(), 10, &types.Message{})
		assert.NoError(t, err)

		var hd = make([]types.Header, 1)
		hd[0] = types.Header{Key: "k", Value: []byte("v")}
		err = p.SendDelay(context.Background(), 10, &types.Message{Headers: hd}, &types.Message{Headers: hd})
		assert.NoError(t, err)
	})
}
