package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const TracerName = "go-zero"

// DoInSpanWithErr executes function doFn inside new span with `operationName` name
// and hooking as child to a span found within given context if any.
// It logs the error inside the new span created, which differentiates it from DoInSpan and DoWithSpan.
func DoInSpanWithErr(ctx context.Context, operationName string, doFn func(context.Context) error,
	opts ...trace.SpanStartOption) error {
	tracer := otel.Tracer(TracerName)
	newCtx, span := tracer.Start(ctx, operationName, opts...)
	defer span.End()

	err := doFn(newCtx)
	if err != nil {
		span.RecordError(err)
	}

	return err
}

// DoInSpan executes function doFn inside new span with `operationName` name
// and hooking as child to a span found within given context if any.
func DoInSpan(ctx context.Context, operationName string, doFn func(context.Context),
	opts ...trace.SpanStartOption) {
	tracer := otel.Tracer(TracerName)
	newCtx, span := tracer.Start(ctx, operationName, opts...)
	defer span.End()
	doFn(newCtx)
}

// DoWithSpan executes function doFn inside new span with `operationName` name
// and hooking as child to a span found within given context if any.
func DoWithSpan(ctx context.Context, operationName string, doFn func(ctx context.Context, span trace.Span),
	opts ...trace.SpanStartOption) {
	tracer := otel.Tracer(TracerName)
	newCtx, span := tracer.Start(ctx, operationName, opts...)
	defer span.End()
	doFn(newCtx, span)
}

// DoWithSpanErr executes function doFn inside new span with `operationName` name
// and hooking as child to a span found within given context if any.
func DoWithSpanErr(ctx context.Context, operationName string,
	doFn func(ctx context.Context, span trace.Span) error,
	opts ...trace.SpanStartOption) error {
	tracer := otel.Tracer(TracerName)
	newCtx, span := tracer.Start(ctx, operationName, opts...)
	defer span.End()

	err := doFn(newCtx, span)
	if err != nil {
		span.RecordError(err)
	}

	return err
}

// DoWithSpanWithTracer executes function doFn inside new span with `operationName` name
// and hooking as child to a span found within given context if any.
func DoWithSpanWithTracer(ctx context.Context, tracer trace.Tracer, operationName string,
	doFn func(ctx context.Context, span trace.Span), opts ...trace.SpanStartOption) {
	newCtx, span := tracer.Start(ctx, operationName, opts...)
	defer span.End()
	doFn(newCtx, span)
}
