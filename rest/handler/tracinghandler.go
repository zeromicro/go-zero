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

// TracingHandler return a middleware that process the opentelemetry.
func TracingHandler(serviceName, path string, traceIgnorePaths []string) func(http.Handler) http.Handler {
	ignorePaths := collection.NewSet()
	ignorePaths.AddStr(traceIgnorePaths...)

	return func(next http.Handler) http.Handler {
		propagator := otel.GetTextMapPropagator()
		tracer := otel.GetTracerProvider().Tracer(trace.TraceName)

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

			trw := &response.WithCodeResponseWriter{Writer: w, Code: http.StatusOK}
			next.ServeHTTP(trw, r.WithContext(spanCtx))

			span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(trw.Code)...)
			span.SetStatus(semconv.SpanStatusFromHTTPStatusCodeAndSpanKind(trw.Code, oteltrace.SpanKindServer))
		})
	}
}
