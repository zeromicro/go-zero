package kafka

import (
	"context"
	"log"
	"runtime/debug"

	zsarama "github.com/zeromicro/go-zero/kafka/internal/sarama"
)

type (
	AsyncProducer interface {
		Producer
		SendWithCallback(ctx context.Context, messages []*Message, handler CallbackHandler) error
	}

	// ProducerError is like sarama.ProducerError,
	// but Msg type is *types.Message.
	ProducerError  = zsarama.ProducerError
	ProducerErrors = zsarama.ProducerErrors
	// SuccessHandler is optional callback function for handling success message
	// from sarama async produce Successes() channel.
	SuccessHandler func(msg *Message)
	// ErrorHandler is optional callback function for handling error message
	// from sarama async produce Errors() channel.
	ErrorHandler func(producerError *ProducerError)
	// CallbackHandler is per message scope callback handler.
	CallbackHandler = zsarama.CallbackHandler
	// AsyncProducerOption is optional option type for AsyncProducer.
	AsyncProducerOption func(a *asyncProducer)

	asyncProducer struct {
		p   *zsarama.AsyncProducer
		pc  ProducerConfig
		opt zsarama.AsyncProducerOption
	}
)

// MustNewAsyncProducer returns an async Producer, exits on any error.
func MustNewAsyncProducer(pc ProducerConfig, opts ...AsyncProducerOption) AsyncProducer {
	p, err := NewAsyncProducer(pc, opts...)
	if err != nil {
		log.Fatalf("%+v\n\n%s", err, debug.Stack())
	}

	return p
}

// NewAsyncProducer create kafka async Producer from config.
func NewAsyncProducer(pc ProducerConfig, opts ...AsyncProducerOption) (AsyncProducer, error) {
	a := &asyncProducer{
		pc: pc,
	}
	for _, opt := range opts {
		opt(a)
	}

	p, err := zsarama.NewAsyncProducer(pc, a.opt)
	if err != nil {
		return nil, err
	}

	a.p = p

	return a, nil
}

// Send sends messages to kafka async producer,
// this API signature is same as sync producer,
// but it always returns nil error immediately,
// you should use WithSuccessHandler and WithErrorHandler to set success/error handler.
func (a *asyncProducer) Send(ctx context.Context, messages ...*Message) error {
	return a.p.Send(ctx, messages...)
}

// SendDelay sends a delay messages to kafka delay queue
func (a *asyncProducer) SendDelay(ctx context.Context, delaySeconds int64, messages ...*Message) error {
	return a.p.SendDelay(ctx, delaySeconds, messages...)
}

// SendWithCallback sends messages to kafka async producer with per message scope callback handler.
func (a *asyncProducer) SendWithCallback(ctx context.Context, messages []*Message, handler CallbackHandler) error {
	return a.p.SendWithCallback(ctx, messages, handler)
}

// WithSuccessHandler sets success handler for async producer.
// Deprecated: use SendWithCallback callback handler instead.
func WithSuccessHandler(h SuccessHandler) AsyncProducerOption {
	return func(a *asyncProducer) {
		a.opt.SuccessHandler = zsarama.SuccessHandler(h)
	}
}

// WithErrorHandler sets error handler for async producer.
// Deprecated: use SendWithCallback callback handler instead.
func WithErrorHandler(h ErrorHandler) AsyncProducerOption {
	return func(a *asyncProducer) {
		a.opt.ErrorHandler = func(producerError *zsarama.ProducerError) {
			h(&ProducerError{
				Err: producerError.Err,
				Msg: producerError.Msg,
			})
		}
	}
}
