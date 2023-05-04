package tracetest

import (
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// NewInMemoryExporter returns a new InMemoryExporter
// and sets it as the global for tests.
func NewInMemoryExporter(t *testing.T) *tracetest.InMemoryExporter {
	me := tracetest.NewInMemoryExporter()
	t.Cleanup(func() {
		me.Reset()
	})
	otel.SetTracerProvider(trace.NewTracerProvider(trace.WithSyncer(me)))

	return me
}
