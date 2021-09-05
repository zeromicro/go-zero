package clientinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	opentelemetry2 "github.com/tal-tech/go-zero/core/trace/opentelemetry"
	"google.golang.org/grpc"
)

func TestOpenTracingInterceptor(t *testing.T) {
	opentelemetry2.StartAgent(opentelemetry2.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})

	cc := new(grpc.ClientConn)
	err := OpenTracingInterceptor()(context.Background(), "/ListUser", nil, nil, cc,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
			opts ...grpc.CallOption) error {
			return nil
		})
	assert.Nil(t, err)
}
