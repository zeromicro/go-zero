package serverinterceptors

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUnaryTimeoutInterceptor(t *testing.T) {
	interceptor := UnaryTimeoutInterceptor(time.Millisecond * 10)
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(
		ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}
