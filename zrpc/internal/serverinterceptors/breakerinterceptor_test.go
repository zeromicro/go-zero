package serverinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStreamBreakerInterceptor(t *testing.T) {
	err := StreamBreakerInterceptor(nil, nil, &grpc.StreamServerInfo{
		FullMethod: "any",
	}, func(_ any, _ grpc.ServerStream) error {
		return status.New(codes.DeadlineExceeded, "any").Err()
	})
	assert.NotNil(t, err)
}

func TestUnaryBreakerInterceptor(t *testing.T) {
	_, err := UnaryBreakerInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "any",
	}, func(_ context.Context, _ any) (any, error) {
		return nil, status.New(codes.DeadlineExceeded, "any").Err()
	})
	assert.NotNil(t, err)
}

func TestUnaryBreakerInterceptorOK(t *testing.T) {
	_, err := UnaryBreakerInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "any",
	}, func(_ context.Context, _ any) (any, error) {
		return nil, nil
	})
	assert.NoError(t, err)
}

func TestUnaryBreakerInterceptor_Unavailable(t *testing.T) {
	_, err := UnaryBreakerInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "any",
	}, func(_ context.Context, _ any) (any, error) {
		return nil, breaker.ErrServiceUnavailable
	})
	assert.NotNil(t, err)
	code := status.Code(err)
	assert.Equal(t, codes.Unavailable, code)
}

func TestUnaryBreakerInterceptor_Deadline(t *testing.T) {
	for i := 0; i < 1000; i++ {
		_, err := UnaryBreakerInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "any",
		}, func(_ context.Context, _ any) (any, error) {
			return nil, context.DeadlineExceeded
		})
		switch status.Code(err) {
		case codes.Unavailable:
		default:
			assert.Equal(t, context.DeadlineExceeded, err)
		}
	}

	var dropped bool
	for i := 0; i < 100; i++ {
		_, err := UnaryBreakerInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "any",
		}, func(_ context.Context, _ any) (any, error) {
			return nil, context.DeadlineExceeded
		})
		switch status.Code(err) {
		case codes.Unavailable:
			dropped = true
		default:
			assert.Equal(t, context.DeadlineExceeded, err)
		}
	}
	assert.True(t, dropped)
}
