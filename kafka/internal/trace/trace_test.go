package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	ztrace "github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/kafka/internal/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

func TestMessageCarrierGet(t *testing.T) {
	testCases := []struct {
		name     string
		carrier  MessageCarrier
		key      string
		expected string
	}{
		{
			name: "exists",
			carrier: NewMessageCarrier(&types.Message{Headers: []types.Header{
				{Key: "foo", Value: []byte("bar")},
			}}),
			key:      "foo",
			expected: "bar",
		},
		{
			name:     "not exists",
			carrier:  NewMessageCarrier(&types.Message{Headers: []types.Header{}}),
			key:      "foo",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.carrier.Get(tc.key)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMessageCarrierSet(t *testing.T) {
	msg := types.Message{Headers: []types.Header{
		{Key: "foo", Value: []byte("bar")},
	}}
	carrier := MessageCarrier{msg: &msg}

	carrier.Set("foo", "bar2")
	carrier.Set("foo2", "bar2")
	carrier.Set("foo2", "bar3")
	carrier.Set("foo3", "bar4")

	assert.ElementsMatch(t, carrier.msg.Headers, []types.Header{
		{Key: "foo", Value: []byte("bar2")},
		{Key: "foo2", Value: []byte("bar3")},
		{Key: "foo3", Value: []byte("bar4")},
	})
}

func TestMessageCarrierKeys(t *testing.T) {
	testCases := []struct {
		name     string
		carrier  MessageCarrier
		expected []string
	}{
		{
			name: "one",
			carrier: MessageCarrier{msg: &types.Message{Headers: []types.Header{
				{Key: "foo", Value: []byte("bar")},
			}}},
			expected: []string{"foo"},
		},
		{
			name:     "none",
			carrier:  MessageCarrier{msg: &types.Message{Headers: []types.Header{}}},
			expected: []string{},
		},
		{
			name: "many",
			carrier: MessageCarrier{msg: &types.Message{Headers: []types.Header{
				{Key: "foo", Value: []byte("bar")},
				{Key: "baz", Value: []byte("quux")},
			}}},
			expected: []string{"foo", "baz"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.carrier.Keys()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_attrsFromMsg(t *testing.T) {
	msg := &types.Message{
		Topic:     "test",
		Partition: 1,
		Offset:    10,
		Key:       []byte("test111"),
	}

	attrs := attrsFromMsg(msg)
	assert.Equal(t, 7, len(attrs))
	for _, attr := range attrs {
		switch attr.Key {
		case semconv.MessagingSystemKey:
			assert.Equal(t, "kafka", attr.Value.AsString())
		case semconv.MessagingDestinationKey:
			assert.Equal(t, "test", attr.Value.AsString())
		case semconv.MessagingMessageIDKey:
			assert.Equal(t, "10", attr.Value.AsString())
		case semconv.MessagingKafkaPartitionKey:
			assert.Equal(t, int64(1), attr.Value.AsInt64())
		case semconv.MessagingKafkaMessageKeyKey:
			assert.Equal(t, "test111", attr.Value.AsString())
		case semconv.MessagingOperationKey:
			assert.Equal(t, "receive", attr.Value.AsString())
		case semconv.MessagingDestinationKindKey:
			assert.Equal(t, "topic", attr.Value.AsString())
		default:
			t.Error("should not reach here")
		}
	}
}

func TestAddConsumerAttrs(t *testing.T) {
	tr := otel.Tracer(ztrace.TracerName)
	_, span := tr.Start(context.Background(), "test")
	AddConsumerAttrs(span, &types.Message{}, attribute.KeyValue{
		Key:   "k1",
		Value: attribute.Value{},
	})
}
