package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/core/opentelemetry"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func OtelHandler(path string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !opentelemetry.Enabled() {
			return next
		}

		propagator := otel.GetTextMapPropagator()
		tracer := otel.GetTracerProvider().Tracer(opentelemetry.TraceName)

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
