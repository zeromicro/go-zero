package handler

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"sync"

	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	traceKeyStatusCode = "http.status_code"
)

var notTracingSpans sync.Map

// DontTraceSpan disable tracing for the specified span name.
func DontTraceSpan(spanName string) {
	notTracingSpans.Store(spanName, lang.Placeholder)
}

type traceResponseWriter struct {
	w    http.ResponseWriter
	code int
}

// Flush implements the http.Flusher interface.
func (w *traceResponseWriter) Flush() {
	if flusher, ok := w.w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Hijack implements the http.Hijacker interface.
func (w *traceResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := w.w.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("server doesn't support hijacking")
}

func (w *traceResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *traceResponseWriter) Write(data []byte) (int, error) {
	return w.w.Write(data)
}

func (w *traceResponseWriter) WriteHeader(statusCode int) {
	w.w.WriteHeader(statusCode)
	w.code = statusCode
}

// TracingHandler return a middleware that process the opentelemetry.
func TracingHandler(serviceName, path string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		propagator := otel.GetTextMapPropagator()
		tracer := otel.GetTracerProvider().Tracer(trace.TraceName)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			spanName := path
			if len(spanName) == 0 {
				spanName = r.URL.Path
			}

			if _, ok := notTracingSpans.Load(spanName); ok {
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
			trw := &traceResponseWriter{
				w:    w,
				code: http.StatusOK,
			}
			next.ServeHTTP(trw, r.WithContext(spanCtx))

			span.SetAttributes(attribute.KeyValue{
				Key:   traceKeyStatusCode,
				Value: attribute.IntValue(trw.code),
			})
			if trw.code >= http.StatusBadRequest {
				span.SetStatus(codes.Error, http.StatusText(trw.code))
			} else {
				span.SetStatus(codes.Ok, "")
			}
		})
	}
}
