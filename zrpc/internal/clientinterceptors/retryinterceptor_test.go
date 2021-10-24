package clientinterceptors

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestRetryInterceptor_WithMax(t *testing.T) {
	n := 4
	for i := 0; i < n; i++ {
		count := 0
		cc := new(grpc.ClientConn)
		err := RetryInterceptor()(context.Background(), "/1", nil, nil, cc,
			func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				count++
				return status.Error(codes.ResourceExhausted, "ResourceExhausted")
			}, RetryWithMax(i))
		assert.Error(t, err)
		assert.Equal(t, i+1, count)
	}

}
func TestRetryInterceptor_Disable(t *testing.T) {
	count := 0
	cc := new(grpc.ClientConn)
	err := RetryInterceptor()(context.Background(), "/1", nil, nil, cc,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			count++
			return status.Error(codes.ResourceExhausted, "ResourceExhausted")
		}, RetryDisable())
	assert.Error(t, err)
	assert.Equal(t, 1, count)
}
