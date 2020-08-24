package internal

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestWithStreamClientInterceptors(t *testing.T) {
	opts := WithStreamClientInterceptors()
	assert.NotNil(t, opts)
}

func TestWithUnaryClientInterceptors(t *testing.T) {
	opts := WithUnaryClientInterceptors()
	assert.NotNil(t, opts)
}

func TestChainStreamClientInterceptors_zero(t *testing.T) {
	interceptors := chainStreamClientInterceptors()
	_, err := interceptors(context.Background(), nil, new(grpc.ClientConn), "/foo",
		func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
			opts ...grpc.CallOption) (grpc.ClientStream, error) {
			return nil, nil
		})
	assert.Nil(t, err)
}

func TestChainStreamClientInterceptors_one(t *testing.T) {
	var called int32
	interceptors := chainStreamClientInterceptors(func(ctx context.Context, desc *grpc.StreamDesc,
		cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (
		grpc.ClientStream, error) {
		atomic.AddInt32(&called, 1)
		return nil, nil
	})
	_, err := interceptors(context.Background(), nil, new(grpc.ClientConn), "/foo",
		func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
			opts ...grpc.CallOption) (grpc.ClientStream, error) {
			return nil, nil
		})
	assert.Nil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&called))
}
