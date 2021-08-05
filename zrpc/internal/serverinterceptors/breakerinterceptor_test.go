package serverinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStreamBreakerInterceptor(t *testing.T) {
	err := StreamBreakerInterceptor(nil, nil, &grpc.StreamServerInfo{
		FullMethod: "any",
	}, func(
		srv interface{}, stream grpc.ServerStream) error {
		return status.New(codes.DeadlineExceeded, "any").Err()
	})
	assert.NotNil(t, err)
}

func TestUnaryBreakerInterceptor(t *testing.T) {
	interceptor := UnaryBreakerInterceptor()
	_, err := interceptor(nil, nil, &grpc.UnaryServerInfo{
		FullMethod: "any",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, status.New(codes.DeadlineExceeded, "any").Err()
	})
	assert.NotNil(t, err)
}
