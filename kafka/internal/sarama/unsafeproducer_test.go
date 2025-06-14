package sarama

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	ztrace "github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/internal/mock/saramamock"
	"github.com/zeromicro/go-zero/kafka/internal/types"
	"go.opentelemetry.io/otel"
)

type mockProducerHandler struct {
	sCount *int
	eCount *int
	d      *chan struct{}
}

func (h *mockProducerHandler) OnProduceSuccess(k string, producedMsg *types.Message) {
	*h.sCount++
	*h.d <- struct{}{}
}
func (h *mockProducerHandler) OnProduceError(k string, producerErr *ProducerError) {
	*h.eCount++
	*h.d <- struct{}{}
}

type testCollector[T any] struct {
	C chan T
}

func newTestCollector[T any](n int) *testCollector[T] {
	return &testCollector[T]{
		C: make(chan T, n),
	}
}

func (t *testCollector[T]) Collect() []T {
	var res []T
	for {
		select {
		case v := <-t.C:
			res = append(res, v)
		default:
			return res
		}
	}
}

func Test_UnsafeProducer_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	mp := saramamock.NewMockAsyncProducer(ctrl)
	tc := newTestCollector[*sarama.ProducerMessage](100)
	mp.EXPECT().Input().Return(tc.C).AnyTimes()
	ap := &AsyncProducer{
		producer:     mp,
		sharedClient: &Client{saramaClient: saramamock.NewMockClient(ctrl)},
	}
	handler0 := &mockProducerHandler{}
	up := UnsafeProducer{
		AsyncProducer: ap,
		k:             "k123",
	}
	msg := &types.Message{
		Topic: "test",
		Key:   sarama.ByteEncoder("abc"),
		Value: sarama.ByteEncoder("abc"),
		Headers: []types.Header{
			{
				Key:   "test",
				Value: []byte("test"),
			},
		},
		Metadata: handler0,
	}
	err := up.Send(context.Background(), msg)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tc.Collect()))
}

func Test_UnsafeProducer_setupCallbacks(t *testing.T) {
	done := make(chan struct{}, 2)
	succssCount := 0
	errCount := 0
	handler0 := &mockProducerHandler{&succssCount, &errCount, &done}

	createMessage := func() *types.Message {
		return &types.Message{
			Topic:     "abc",
			Key:       sarama.ByteEncoder("xxx"),
			Value:     sarama.ByteEncoder("xxx"),
			Offset:    100,
			Partition: 0,
			Metadata:  handler0,
		}
	}

	t.Run("with global callback", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		sc := saramamock.NewMockAsyncProducer(ctrl)
		successCh := make(chan *sarama.ProducerMessage, 1)
		errorCh := make(chan *sarama.ProducerError, 1)
		sc.EXPECT().Errors().Return(errorCh).AnyTimes()
		sc.EXPECT().Successes().Return(successCh).AnyTimes()
		c := &AsyncProducer{
			producer:     sc,
			closeCh:      make(chan struct{}),
			sharedClient: &Client{saramaClient: saramamock.NewMockClient(ctrl)},
		}
		uc := &UnsafeProducer{
			AsyncProducer: c,
			k:             "k123",
		}
		tracer := otel.Tracer(ztrace.TracerName)

		uc.setupCallbacks()
		sMsg := message2ProducerMessage(createMessage())
		_, span := tracer.Start(context.Background(), fmt.Sprintf("%s publish", sMsg.Topic))
		setInternalMetadata(sMsg, &producerMetadata{
			startTime: time.Now(),
			callbackHandler: func(msg *types.Message, err error) {
				if msg.Metadata != nil {
					pcb, _ := msg.Metadata.(ProduceHandler)
					if err != nil {
						pcb.OnProduceError(uc.k, &ProducerError{
							Msg: msg,
							Err: err,
						})
					} else {
						pcb.OnProduceSuccess(uc.k, msg)
					}
				}
			},
			span: span,
		})

		successCh <- sMsg
		errorCh <- &sarama.ProducerError{
			Msg: sMsg,
			Err: errors.New("test"),
		}
		<-done
		<-done
		assert.Equal(t, 1, succssCount)
		assert.Equal(t, 1, errCount)
	})
}

func TestUnsafeProducer_SendDelay(t *testing.T) {
	ctrl := gomock.NewController(t)
	ap := saramamock.NewMockAsyncProducer(ctrl)
	tc := newTestCollector[*sarama.ProducerMessage](100)
	ap.EXPECT().Input().Return(tc.C).AnyTimes()
	p := &AsyncProducer{
		producer:     ap,
		sharedClient: &Client{saramaClient: saramamock.NewMockClient(ctrl)},
	}

	up := &UnsafeProducer{
		p,
		"k123",
	}

	t.Run("send fail", func(t *testing.T) {
		err := up.SendDelay(context.Background(), 3, &types.Message{})
		assert.Error(t, err)
	})
	t.Run("send success", func(t *testing.T) {
		err := up.SendDelay(context.Background(), 10, &types.Message{})
		assert.NoError(t, err)
	})
}
