package clientinterceptors

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestTracingInterceptor(t *testing.T) {
	var run int32
	var wg sync.WaitGroup
	wg.Add(1)
	cc := new(grpc.ClientConn)
	err := TracingInterceptor(context.Background(), "/foo", nil, nil, cc,
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

func TestTracingInterceptor_GrpcFormat(t *testing.T) {
	var run int32
	var wg sync.WaitGroup
	wg.Add(1)
	md := metadata.New(map[string]string{
		"foo": "bar",
	})
	carrier, err := trace.Inject(trace.GrpcFormat, md)
	assert.Nil(t, err)
	ctx, _ := trace.StartServerSpan(context.Background(), carrier, "user", "/foo")
	cc := new(grpc.ClientConn)
	err = TracingInterceptor(ctx, "/foo", nil, nil, cc,
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
