package trace

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestSpanIDFromContext(t *testing.T) {
	tracer := sdktrace.NewTracerProvider().Tracer("test")
	ctx, span := tracer.Start(
		context.Background(),
		"foo",
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
		oteltrace.WithAttributes(semconv.HTTPClientAttributesFromHTTPRequest(httptest.NewRequest(http.MethodGet, "/", nil))...),
	)
	defer span.End()

	assert.NotEmpty(t, TraceIDFromContext(ctx))
	assert.NotEmpty(t, SpanIDFromContext(ctx))
}

func TestSpanIDFromContextEmpty(t *testing.T) {
	assert.Empty(t, TraceIDFromContext(context.Background()))
	assert.Empty(t, SpanIDFromContext(context.Background()))
}
