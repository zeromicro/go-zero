package handler

import (
	"net/http"

	opentelemetry2 "github.com/tal-tech/go-zero/core/trace/opentelemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// OtelHandler return a middleware that process the opentelemetry.
func OtelHandler(path string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !opentelemetry2.Enabled() {
			return next
		}

		propagator := otel.GetTextMapPropagator()
		tracer := otel.GetTracerProvider().Tracer(opentelemetry2.TraceName)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			spanCtx, span := tracer.Start(
				ctx,
				path,
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
				oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest("", path, r)...),
			)
			defer span.End()

			next.ServeHTTP(w, r.WithContext(spanCtx))
		})
	}
}
