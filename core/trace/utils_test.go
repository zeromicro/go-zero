package trace

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/peer"
)

func TestPeerFromContext(t *testing.T) {
	addrs, err := net.InterfaceAddrs()
	assert.Nil(t, err)
	assert.NotEmpty(t, addrs)
	tests := []struct {
		name  string
		ctx   context.Context
		empty bool
	}{
		{
			name:  "empty",
			ctx:   context.Background(),
			empty: true,
		},
		{
			name:  "nil",
			ctx:   peer.NewContext(context.Background(), nil),
			empty: true,
		},
		{
			name: "with value",
			ctx: peer.NewContext(context.Background(), &peer.Peer{
				Addr: addrs[0],
			}),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			addr := PeerFromCtx(test.ctx)
			assert.Equal(t, test.empty, len(addr) == 0)
		})
	}
}

func TestParseFullMethod(t *testing.T) {
	tests := []struct {
		fullMethod string
		name       string
		attr       []attribute.KeyValue
	}{
		{
			fullMethod: "/grpc.test.EchoService/Echo",
			name:       "grpc.test.EchoService/Echo",
			attr: []attribute.KeyValue{
				semconv.RPCServiceKey.String("grpc.test.EchoService"),
				semconv.RPCMethodKey.String("Echo"),
			},
		}, {
			fullMethod: "/com.example.ExampleRmiService/exampleMethod",
			name:       "com.example.ExampleRmiService/exampleMethod",
			attr: []attribute.KeyValue{
				semconv.RPCServiceKey.String("com.example.ExampleRmiService"),
				semconv.RPCMethodKey.String("exampleMethod"),
			},
		}, {
			fullMethod: "/MyCalcService.Calculator/Add",
			name:       "MyCalcService.Calculator/Add",
			attr: []attribute.KeyValue{
				semconv.RPCServiceKey.String("MyCalcService.Calculator"),
				semconv.RPCMethodKey.String("Add"),
			},
		}, {
			fullMethod: "/MyServiceReference.ICalculator/Add",
			name:       "MyServiceReference.ICalculator/Add",
			attr: []attribute.KeyValue{
				semconv.RPCServiceKey.String("MyServiceReference.ICalculator"),
				semconv.RPCMethodKey.String("Add"),
			},
		}, {
			fullMethod: "/MyServiceWithNoPackage/theMethod",
			name:       "MyServiceWithNoPackage/theMethod",
			attr: []attribute.KeyValue{
				semconv.RPCServiceKey.String("MyServiceWithNoPackage"),
				semconv.RPCMethodKey.String("theMethod"),
			},
		}, {
			fullMethod: "/pkg.svr",
			name:       "pkg.svr",
			attr:       []attribute.KeyValue(nil),
		}, {
			fullMethod: "/pkg.svr/",
			name:       "pkg.svr/",
			attr: []attribute.KeyValue{
				semconv.RPCServiceKey.String("pkg.svr"),
			},
		},
	}

	for _, test := range tests {
		n, a := ParseFullMethod(test.fullMethod)
		assert.Equal(t, test.name, n)
		assert.Equal(t, test.attr, a)
	}
}

func TestSpanInfo(t *testing.T) {
	val, kvs := SpanInfo("/fullMethod", "remote")
	assert.Equal(t, "fullMethod", val)
	assert.NotEmpty(t, kvs)
}

func TestPeerAttr(t *testing.T) {
	tests := []struct {
		name   string
		addr   string
		expect []attribute.KeyValue
	}{
		{
			name: "empty",
		},
		{
			name: "port only",
			addr: ":8080",
			expect: []attribute.KeyValue{
				semconv.NetPeerIPKey.String(localhost),
				semconv.NetPeerPortKey.String("8080"),
			},
		},
		{
			name: "port only",
			addr: "192.168.0.2:8080",
			expect: []attribute.KeyValue{
				semconv.NetPeerIPKey.String("192.168.0.2"),
				semconv.NetPeerPortKey.String("8080"),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			kvs := PeerAttr(test.addr)
			assert.EqualValues(t, test.expect, kvs)
		})
	}
}

func TestTracerFromContext(t *testing.T) {
	traceFn := func(ctx context.Context, hasTraceId bool) {
		spanContext := trace.SpanContextFromContext(ctx)
		assert.Equal(t, spanContext.IsValid(), hasTraceId)
		parentTraceId := spanContext.TraceID().String()

		tracer := TracerFromContext(ctx)
		_, span := tracer.Start(ctx, "b")
		defer span.End()

		spanContext = span.SpanContext()
		assert.True(t, spanContext.IsValid())
		if hasTraceId {
			assert.Equal(t, parentTraceId, spanContext.TraceID().String())
		}

	}

	t.Run("context", func(t *testing.T) {
		opts := []sdktrace.TracerProviderOption{
			// Set the sampling rate based on the parent span to 100%
			sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1))),
			// Record information about this application in a Resource.
			sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String("test"))),
		}
		tp = sdktrace.NewTracerProvider(opts...)
		otel.SetTracerProvider(tp)
		ctx, span := tp.Tracer(TraceName).Start(context.Background(), "a")

		defer span.End()
		traceFn(ctx, true)
	})

	t.Run("global", func(t *testing.T) {
		opts := []sdktrace.TracerProviderOption{
			// Set the sampling rate based on the parent span to 100%
			sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1))),
			// Record information about this application in a Resource.
			sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String("test"))),
		}
		tp = sdktrace.NewTracerProvider(opts...)
		otel.SetTracerProvider(tp)

		traceFn(context.Background(), false)
	})
}
