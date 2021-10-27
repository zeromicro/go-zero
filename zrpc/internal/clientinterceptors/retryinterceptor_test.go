package clientinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRetryInterceptor_WithMax(t *testing.T) {
	n := 4
	for i := 0; i < n; i++ {
		count := 0
		cc := new(grpc.ClientConn)
		err := RetryInterceptor(true)(context.Background(), "/1", nil, nil, cc,
			func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
				opts ...grpc.CallOption) error {
				count++
				return status.Error(codes.ResourceExhausted, "ResourceExhausted")
			}, retry.WithMax(i))
		assert.Error(t, err)
		assert.Equal(t, i+1, count)
	}
}
