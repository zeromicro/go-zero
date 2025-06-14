package trace

import (
	"strconv"

	"github.com/zeromicro/go-zero/kafka/internal/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// copy from https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/Shopify/sarama/otelsarama/message.go

var _ propagation.TextMapCarrier = (*MessageCarrier)(nil)

// MessageCarrier injects and extracts traces from a types.Message.
type MessageCarrier struct {
	msg *types.Message
}

// NewMessageCarrier creates a new MessageCarrier.
func NewMessageCarrier(msg *types.Message) MessageCarrier {
	return MessageCarrier{msg: msg}
}

// Get retrieves a single value for a given key.
func (c MessageCarrier) Get(key string) string {
	return c.msg.GetHeader(key)
}

// Set sets a header.
func (c MessageCarrier) Set(key, val string) {
	c.msg.SetHeader(key, val)
}

// Keys returns a slice of all key identifiers in the carrier.
func (c MessageCarrier) Keys() []string {
	out := make([]string, len(c.msg.Headers))
	for i, h := range c.msg.Headers {
		out[i] = h.Key
	}
	return out
}

func attrsFromMsg(msg *types.Message) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.MessagingSystemKey.String("kafka"),
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationKey.String(msg.Topic),
		semconv.MessagingOperationReceive,
		semconv.MessagingMessageIDKey.String(strconv.FormatInt(msg.Offset, 10)),
		semconv.MessagingKafkaPartitionKey.Int64(int64(msg.Partition)),
		semconv.MessagingKafkaMessageKeyKey.String(string(msg.Key)),
	}
}

func AddConsumerAttrs(span trace.Span, msg *types.Message, extAttrs ...attribute.KeyValue) {
	attrs := append(extAttrs,
		attrsFromMsg(msg)...,
	)

	span.SetAttributes(attrs...)
}
