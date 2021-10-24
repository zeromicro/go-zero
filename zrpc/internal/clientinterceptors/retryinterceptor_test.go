package clientinterceptors

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestRetryInterceptor(t *testing.T) {
	cc := new(grpc.ClientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := RetryInterceptor()(ctx, "/1", nil, nil, cc,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			time.Sleep(time.Millisecond)
			return status.Error(codes.ResourceExhausted, "ResourceExhausted")
		}, WithMax(1000), WithPerRetryTimeout(time.Microsecond))
	assert.Error(t, err)

}
