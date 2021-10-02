package trace

import (
	"fmt"
	"sync"

	"github.com/tal-tech/go-zero/core/logx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const kindJaeger = "jaeger"

var once sync.Once

// StartAgent starts a opentelemetry agent.
func StartAgent(c Config) {
	once.Do(func() {
		startAgent(c)
	})
}

func createExporter(c Config) (sdktrace.SpanExporter, error) {
	// Just support jaeger now, more for later
	switch c.Batcher {
	case kindJaeger:
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(c.Endpoint)))
		if err != nil {
			return nil, err
		}

		return exp, nil
	default:
		return nil, fmt.Errorf("unknown exporter: %s", c.Batcher)
	}
}

func startAgent(c Config) {
	opts := []sdktrace.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(c.Sampler))),
		// Record information about this application in an Resource.
		sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(c.Name))),
	}

	if len(c.Endpoint) > 0 {
		exp, err := createExporter(c)
		if err != nil {
			logx.Error(err)
			return
		}

		opts = append(opts,
			// Always be sure to batch in production.
			sdktrace.WithBatcher(exp),
		)
	}

	tp := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logx.Errorf("[otel] error: %v", err)
	}))
}
