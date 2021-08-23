package serverinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/opentelemetry"
	"google.golang.org/grpc"
)

func TestUnaryOpenTracingInterceptor_Disable(t *testing.T) {
	interceptor := UnaryOpenTracingInterceptor()
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestUnaryOpenTracingInterceptor_Enabled(t *testing.T) {
	opentelemetry.StartAgent(opentelemetry.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})
	interceptor := UnaryOpenTracingInterceptor()
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/package.TestService.GetUser",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}
