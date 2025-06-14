package sarama

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/zeromicro/go-zero/internal/mock/saramamock"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zeromicro/go-zero/kafka/internal/types"
)

func TestProducerError(t *testing.T) {
	e := errors.New("test error")
	err := ProducerError{
		Err: e,
		Msg: &types.Message{
			Topic: "test-topic",
		},
	}
	assert.Contains(t, err.Error(), "test error")
	assert.Contains(t, err.Error(), "test-topic")
	assert.Equal(t, e, errors.Unwrap(err))
}

func TestProducerErrors(t *testing.T) {
	var errs ProducerErrors
	errs = append(errs, &ProducerError{
		Err: errors.New("test error 1"),
		Msg: &types.Message{
			Topic: "test-topic-1",
		},
	}, &ProducerError{
		Err: errors.New("test error 2"),
		Msg: &types.Message{
			Topic: "test-topic-2",
		},
	})
	assert.Contains(t, errs.Error(), "deliver 2 messages")
}

func TestNewAsyncProducer(t *testing.T) {
	t.Run("no panic", func(t *testing.T) {
		_, err := NewAsyncProducer(types.ProducerConfig{
			Client:               types.ClientConfig{Brokers: []string{}},
			SharedProducerConfig: types.SharedProducerConfig{EnableRecovery: false},
		}, AsyncProducerOption{
			SuccessHandler: func(msg *types.Message) {
				panic(errors.New("mock"))
			},
		})
		assert.Error(t, err)
	})
}

func TestNewProducer_SendDelay(t *testing.T) {
	ctrl := gomock.NewController(t)
	ap := saramamock.NewMockAsyncProducer(ctrl)
	tc := newTestCollector[*sarama.ProducerMessage](100)
	ap.EXPECT().Input().Return(tc.C).AnyTimes()
	p := &AsyncProducer{
		producer:     ap,
		sharedClient: &Client{saramaClient: saramamock.NewMockClient(ctrl)},
	}

	t.Run("send fail", func(t *testing.T) {
		err := p.SendDelay(context.Background(), 3)
		assert.Error(t, err)

		err = p.SendDelay(context.Background(), 3, &types.Message{})
		assert.Error(t, err)
	})
	t.Run("send success", func(t *testing.T) {
		err := p.SendDelay(context.Background(), 10, &types.Message{})
		assert.NoError(t, err)
		var hd = make([]types.Header, 1)
		hd[0] = types.Header{Key: "k", Value: []byte("v")}
		err = p.SendDelay(context.Background(), 10, &types.Message{Headers: hd}, &types.Message{Headers: hd})
		assert.NoError(t, err)
	})
}
