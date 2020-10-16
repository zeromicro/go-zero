package serverinterceptors

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/trace/tracespec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnaryTracingInterceptor(t *testing.T) {
	interceptor := UnaryTracingInterceptor("foo")
	var run int32
	var wg sync.WaitGroup
	wg.Add(1)
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		defer wg.Done()
		atomic.AddInt32(&run, 1)
		return nil, nil
	})
	wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&run))
}

func TestUnaryTracingInterceptor_GrpcFormat(t *testing.T) {
	interceptor := UnaryTracingInterceptor("foo")
	var wg sync.WaitGroup
	wg.Add(1)
	var md metadata.MD
	ctx := metadata.NewIncomingContext(context.Background(), md)
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		defer wg.Done()
		assert.True(t, len(ctx.Value(tracespec.TracingKey).(tracespec.Trace).TraceId()) > 0)
		assert.True(t, len(ctx.Value(tracespec.TracingKey).(tracespec.Trace).SpanId()) > 0)
		return nil, nil
	})
	wg.Wait()
	assert.Nil(t, err)
}
