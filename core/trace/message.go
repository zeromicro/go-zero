package trace

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

const messageEvent = "message"

var (
	// MessageSent is the type of sent messages.
	MessageSent = messageType(RPCMessageTypeSent)
	// MessageReceived is the type of received messages.
	MessageReceived = messageType(RPCMessageTypeReceived)
)

type messageType attribute.KeyValue

// Event adds an event of the messageType to the span associated with the
// passed context with id and size (if message is a proto message).
func (m messageType) Event(ctx context.Context, id int, message any) {
	span := trace.SpanFromContext(ctx)
	if p, ok := message.(proto.Message); ok {
		span.AddEvent(messageEvent, trace.WithAttributes(
			attribute.KeyValue(m),
			RPCMessageIDKey.Int(id),
			RPCMessageUncompressedSizeKey.Int(proto.Size(p)),
		))
	} else {
		span.AddEvent(messageEvent, trace.WithAttributes(
			attribute.KeyValue(m),
			RPCMessageIDKey.Int(id),
		))
	}
}
