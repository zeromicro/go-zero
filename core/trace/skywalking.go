package trace

import (
	"context"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	"github.com/SkyAPM/go2sky/reporter"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Exporter struct {
	Tracer *go2sky.Tracer
}

var _ sdktrace.SpanExporter = (*Exporter)(nil)

func NewSkywalking(endpoint, serviceName string) (*Exporter, error) {
	rp, err := reporter.NewGRPCReporter(endpoint, reporter.WithCheckInterval(time.Second))
	if err != nil {
		return nil, err
	}
	tracer, err := go2sky.NewTracer(serviceName, go2sky.WithReporter(rp))
	e := &Exporter{
		Tracer: tracer,
	}
	return e, nil
}

func (e *Exporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	for _, s := range spans {
		span, _, err := e.Tracer.CreateEntrySpan(ctx, s.Name(), func(key string) (string, error) {
			scx := propagation.SpanContext{}
			if !s.Parent().TraceID().IsValid() { //parent
				spanid := ([8]byte)(s.SpanContext().SpanID())
				sid := crc32.ChecksumIEEE(spanid[:]) / 2
				scx = propagation.SpanContext{
					Sample:                1,
					TraceID:               s.SpanContext().TraceID().String(),
					ParentSegmentID:       s.Parent().SpanID().String(),
					ParentSpanID:          int32(sid),
					ParentService:         s.Name(),
					ParentServiceInstance: s.Name(),
					ParentEndpoint:        s.Name(),
					AddressUsedAtClient:   s.Name(),
				}
			} else { //child
				spanid := ([8]byte)(s.Parent().SpanID())
				sid := crc32.ChecksumIEEE(spanid[:]) / 2
				scx = propagation.SpanContext{
					Sample:                1,
					TraceID:               s.SpanContext().TraceID().String(),
					ParentSegmentID:       s.Parent().SpanID().String(),
					ParentSpanID:          int32(sid),
					ParentService:         s.Name(),
					ParentServiceInstance: s.Name(),
					ParentEndpoint:        s.Name(),
					AddressUsedAtClient:   s.Name(),
				}
			}

			return scx.EncodeSW8(), nil
		})
		if err != nil {
			fmt.Println("err:", err)
		}
		span.SetComponent(8888)
		span.Tag(go2sky.TagURL, s.Name())
		span.SetSpanLayer(0)
		span.Tag(go2sky.TagStatusCode, "200")
		span.Tag(go2sky.TagURL, s.Name())
		span.End()
	}
	return nil
}

// Shutdown stops the exporter flushing any pending exports.
func (e *Exporter) Shutdown(ctx context.Context) error {
	return nil
}
