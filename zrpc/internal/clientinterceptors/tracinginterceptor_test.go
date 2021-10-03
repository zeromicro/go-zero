package clientinterceptors

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/trace"
	"google.golang.org/grpc"
)

func TestOpenTracingInterceptor(t *testing.T) {
	trace.StartAgent(trace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})

	cc := new(grpc.ClientConn)
	err := UnaryTracingInterceptor(context.Background(), "/ListUser", nil, nil, cc,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
			opts ...grpc.CallOption) error {
			return nil
		})
	assert.Nil(t, err)
}

func TestUnaryTracingInterceptor(t *testing.T) {
	var run int32
	var wg sync.WaitGroup
	wg.Add(1)
	cc := new(grpc.ClientConn)
	err := UnaryTracingInterceptor(context.Background(), "/foo", nil, nil, cc,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
			opts ...grpc.CallOption) error {
			defer wg.Done()
			atomic.AddInt32(&run, 1)
			return nil
		})
	wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&run))
}

func TestStreamTracingInterceptor(t *testing.T) {
	var run int32
	var wg sync.WaitGroup
	wg.Add(1)
	cc := new(grpc.ClientConn)
	_, err := StreamTracingInterceptor(context.Background(), nil, cc, "/foo",
		func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
			opts ...grpc.CallOption) (grpc.ClientStream, error) {
			defer wg.Done()
			atomic.AddInt32(&run, 1)
			return nil, nil
		})
	wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&run))
}

func TestUnaryTracingInterceptor_GrpcFormat(t *testing.T) {
	var run int32
	var wg sync.WaitGroup
	wg.Add(1)
	cc := new(grpc.ClientConn)
	err := UnaryTracingInterceptor(context.Background(), "/foo", nil, nil, cc,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
			opts ...grpc.CallOption) error {
			defer wg.Done()
			atomic.AddInt32(&run, 1)
			return nil
		})
	wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&run))
}

func TestStreamTracingInterceptor_GrpcFormat(t *testing.T) {
	var run int32
	var wg sync.WaitGroup
	wg.Add(1)
	cc := new(grpc.ClientConn)
	_, err := StreamTracingInterceptor(context.Background(), nil, cc, "/foo",
		func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
			opts ...grpc.CallOption) (grpc.ClientStream, error) {
			defer wg.Done()
			atomic.AddInt32(&run, 1)
			return nil, nil
		})
	wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&run))
}
