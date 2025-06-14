package sarama

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/kafka/internal/types"
	otrace "go.opentelemetry.io/otel/trace"
)

type (
	ProduceHandler interface {
		OnProduceSuccess(k string, producedMsg *types.Message)
		OnProduceError(k string, producerErr *ProducerError)
	}
	// UnsafeProducer wrapper sarama async producer.
	UnsafeProducer struct {
		*AsyncProducer
		k string
	}
)

func newUnsafeProducerFromClient(c *Client, k string, opt AsyncProducerOption) (*UnsafeProducer, error) {
	cc, err := newAsyncProducerFromClient(c, k, opt)
	if err != nil {
		return nil, err
	}
	ecc := &UnsafeProducer{
		AsyncProducer: cc,
		k:             k,
	}
	return ecc, nil
}

func (a *UnsafeProducer) Send(ctx context.Context, messages ...*types.Message) error {
	return a.sendWithCallback(ctx, messages, func(sMsg *sarama.ProducerMessage, message *types.Message, span otrace.Span) {
		setInternalMetadata(sMsg, &producerMetadata{
			startTime: time.Now(),
			callbackHandler: func(msg *types.Message, err error) {
				if msg.Metadata != nil {
					pcb, _ := msg.Metadata.(ProduceHandler)
					if err != nil {
						pcb.OnProduceError(a.k, &ProducerError{
							Msg: msg,
							Err: err,
						})
					} else {
						pcb.OnProduceSuccess(a.k, msg)
					}
				}
			},
			span: span,
		})
	})
}

// setupCallbacks setup callbacks loop for async producer.
func (a *UnsafeProducer) setupCallbacks() {
	a.setupCallbacksRunSafe(true)
}
