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

func TestStreamTracingInterceptor(t *testing.T) {
	interceptor := StreamTracingInterceptor("foo")
	var run int32
	var wg sync.WaitGroup
	wg.Add(1)
	err := interceptor(nil, new(mockedServerStream), nil,
		func(srv interface{}, stream grpc.ServerStream) error {
			defer wg.Done()
			atomic.AddInt32(&run, 1)
			return nil
		})
	wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&run))
}

func TestStreamTracingInterceptor_GrpcFormat(t *testing.T) {
	interceptor := StreamTracingInterceptor("foo")
	var run int32
	var wg sync.WaitGroup
	wg.Add(1)
	var md metadata.MD
	ctx := metadata.NewIncomingContext(context.Background(), md)
	stream := mockedServerStream{ctx: ctx}
	err := interceptor(nil, &stream, &grpc.StreamServerInfo{
		FullMethod: "/foo",
	}, func(srv interface{}, stream grpc.ServerStream) error {
		defer wg.Done()
		atomic.AddInt32(&run, 1)
		return nil
	})
	wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&run))
}

type mockedServerStream struct {
	ctx context.Context
}

func (m *mockedServerStream) SetHeader(md metadata.MD) error {
	panic("implement me")
}

func (m *mockedServerStream) SendHeader(md metadata.MD) error {
	panic("implement me")
}

func (m *mockedServerStream) SetTrailer(md metadata.MD) {
	panic("implement me")
}

func (m *mockedServerStream) Context() context.Context {
	if m.ctx == nil {
		return context.Background()
	}

	return m.ctx
}

func (m *mockedServerStream) SendMsg(v interface{}) error {
	panic("implement me")
}

func (m *mockedServerStream) RecvMsg(v interface{}) error {
	panic("implement me")
}
