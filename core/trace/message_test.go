package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func TestMessageType_Event(t *testing.T) {
	ctx, s := otel.Tracer(TraceName).Start(context.Background(), "test")
	span := mockSpan{Span: s}
	ctx = trace.ContextWithSpan(ctx, &span)
	MessageReceived.Event(ctx, 1, "foo")
	assert.Equal(t, messageEvent, span.name)
	assert.NotEmpty(t, span.options)
}

func TestMessageType_EventProtoMessage(t *testing.T) {
	var span mockSpan
	var message mockMessage
	ctx := trace.ContextWithSpan(context.Background(), &span)
	MessageReceived.Event(ctx, 1, message)
	assert.Equal(t, messageEvent, span.name)
	assert.NotEmpty(t, span.options)
}

type mockSpan struct {
	trace.Span
	name    string
	options []trace.EventOption
}

func (m *mockSpan) End(_ ...trace.SpanEndOption) {
}

func (m *mockSpan) AddEvent(name string, options ...trace.EventOption) {
	m.name = name
	m.options = options
}

func (m *mockSpan) IsRecording() bool {
	return false
}

func (m *mockSpan) RecordError(_ error, _ ...trace.EventOption) {
}

func (m *mockSpan) SpanContext() trace.SpanContext {
	panic("implement me")
}

func (m *mockSpan) SetStatus(_ codes.Code, _ string) {
}

func (m *mockSpan) SetName(_ string) {
}

func (m *mockSpan) SetAttributes(_ ...attribute.KeyValue) {
}

func (m *mockSpan) TracerProvider() trace.TracerProvider {
	return nil
}

type mockMessage struct{}

func (m mockMessage) ProtoReflect() protoreflect.Message {
	return new(dynamicpb.Message)
}
