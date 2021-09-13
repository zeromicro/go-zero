package opentelemetry

import (
	"sync"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/syncx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	once    sync.Once
	enabled syncx.AtomicBool
)

// Enabled returns if opentelemetry is enabled.
func Enabled() bool {
	return enabled.True()
}

// StartAgent starts a opentelemetry agent.
func StartAgent(c Config) {
	once.Do(func() {
		if len(c.Endpoint) == 0 {
			return
		}

		// Just support jaeger now
		if c.Batcher != "jaeger" {
			return
		}

		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(c.Endpoint)))
		if err != nil {
			logx.Error(err)
			return
		}

		tp := sdktrace.NewTracerProvider(
			// Set the sampling rate based on the parent span to 100%
			sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(c.Sampler))),
			// Always be sure to batch in production.
			sdktrace.WithBatcher(exp),
			// Record information about this application in an Resource.
			sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(c.Name))),
		)

		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
		otel.SetErrorHandler(otelErrHandler{})

		enabled.Set(true)
	})
}

// errHandler handing otel errors.
type otelErrHandler struct{}

var _ otel.ErrorHandler = otelErrHandler{}

func (o otelErrHandler) Handle(err error) {
	logx.Errorf("[otel] error: %v", err)
}
