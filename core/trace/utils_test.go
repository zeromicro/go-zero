package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

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
			fullMethod: "/pkg.srv",
			name:       "pkg.srv",
			attr:       []attribute.KeyValue(nil),
		}, {
			fullMethod: "/pkg.srv/",
			name:       "pkg.srv/",
			attr: []attribute.KeyValue{
				semconv.RPCServiceKey.String("pkg.srv"),
			},
		},
	}

	for _, test := range tests {
		n, a := ParseFullMethod(test.fullMethod)
		assert.Equal(t, test.name, n)
		assert.Equal(t, test.attr, a)
	}
}
