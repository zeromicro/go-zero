package serverinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stat"
	"google.golang.org/grpc"
)

func TestUnaryStatInterceptor(t *testing.T) {
	metrics := stat.NewMetrics("mock")
	interceptor := UnaryStatInterceptor(metrics)
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestUnaryStatInterceptor_crash(t *testing.T) {
	metrics := stat.NewMetrics("mock")
	interceptor := UnaryStatInterceptor(metrics)
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("error")
	})
	assert.NotNil(t, err)
}
