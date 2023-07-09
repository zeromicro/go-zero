package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/rest/internal/response"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type (
	// TraceOption defines the method to customize an traceOptions.
	TraceOption func(options *traceOptions)

	// traceOptions is TraceHandler options.
	traceOptions struct {
		traceIgnorePaths []string
	}
)

// TraceHandler return a middleware that process the opentelemetry.
func TraceHandler(serviceName, path string, opts ...TraceOption) func(http.Handler) http.Handler {
	var options traceOptions
	for _, opt := range opts {
		opt(&options)
	}

	ignorePaths := collection.NewSet()
	ignorePaths.AddStr(options.traceIgnorePaths...)

	return func(next http.Handler) http.Handler {
		tracer := otel.Tracer(trace.TraceName)
		propagator := otel.GetTextMapPropagator()

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			spanName := path
			if len(spanName) == 0 {
				spanName = r.URL.Path
			}

			if ignorePaths.Contains(spanName) {
				next.ServeHTTP(w, r)
				return
			}

			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			spanCtx, span := tracer.Start(
				ctx,
				spanName,
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
				oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(
					serviceName, spanName, r)...),
			)
			defer span.End()

			// convenient for tracking error messages
			propagator.Inject(spanCtx, propagation.HeaderCarrier(w.Header()))

			trw := response.NewWithCodeResponseWriter(w)
			next.ServeHTTP(trw, r.WithContext(spanCtx))

			span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(trw.Code)...)
			span.SetStatus(semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(
				trw.Code, oteltrace.SpanKindServer))
		})
	}
}

// WithTraceIgnorePaths specifies the traceIgnorePaths option for TraceHandler.
func WithTraceIgnorePaths(traceIgnorePaths []string) TraceOption {
	return func(options *traceOptions) {
		options.traceIgnorePaths = append(options.traceIgnorePaths, traceIgnorePaths...)
	}
}
