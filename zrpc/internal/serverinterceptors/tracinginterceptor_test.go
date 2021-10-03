package serverinterceptors

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

func TestUnaryOpenTracingInterceptor_Disable(t *testing.T) {
	_, err := UnaryTracingInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestUnaryOpenTracingInterceptor_Enabled(t *testing.T) {
	trace.StartAgent(trace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})
	_, err := UnaryTracingInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/package.TestService.GetUser",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestUnaryTracingInterceptor(t *testing.T) {
	var run int32
	var wg sync.WaitGroup
	wg.Add(1)
	_, err := UnaryTracingInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
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

func TestStreamTracingInterceptor_GrpcFormat(t *testing.T) {
	var run int32
	var wg sync.WaitGroup
	wg.Add(1)
	var md metadata.MD
	ctx := metadata.NewIncomingContext(context.Background(), md)
	stream := mockedServerStream{ctx: ctx}
	err := StreamTracingInterceptor(nil, &stream, &grpc.StreamServerInfo{
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
