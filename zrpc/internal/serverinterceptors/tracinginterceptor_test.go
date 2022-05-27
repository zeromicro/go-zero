package serverinterceptors

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func TestUnaryTracingInterceptor_WithError(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "normal error",
			err:  errors.New("dummy"),
		},
		{
			name: "grpc error",
			err:  status.Error(codes.DataLoss, "dummy"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var wg sync.WaitGroup
			wg.Add(1)
			var md metadata.MD
			ctx := metadata.NewIncomingContext(context.Background(), md)
			_, err := UnaryTracingInterceptor(ctx, nil, &grpc.UnaryServerInfo{
				FullMethod: "/",
			}, func(ctx context.Context, req interface{}) (interface{}, error) {
				defer wg.Done()
				return nil, test.err
			})
			wg.Wait()
			assert.Equal(t, test.err, err)
		})
	}
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
	}, func(svr interface{}, stream grpc.ServerStream) error {
		defer wg.Done()
		atomic.AddInt32(&run, 1)
		return nil
	})
	wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&run))
}

func TestStreamTracingInterceptor_FinishWithGrpcError(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "receive event",
			err:  status.Error(codes.DataLoss, "dummy"),
		},
		{
			name: "error event",
			err:  status.Error(codes.DataLoss, "dummy"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var wg sync.WaitGroup
			wg.Add(1)
			var md metadata.MD
			ctx := metadata.NewIncomingContext(context.Background(), md)
			stream := mockedServerStream{ctx: ctx}
			err := StreamTracingInterceptor(nil, &stream, &grpc.StreamServerInfo{
				FullMethod: "/foo",
			}, func(svr interface{}, stream grpc.ServerStream) error {
				defer wg.Done()
				return test.err
			})
			wg.Wait()
			assert.Equal(t, test.err, err)
		})
	}
}

func TestStreamTracingInterceptor_WithError(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "normal error",
			err:  errors.New("dummy"),
		},
		{
			name: "grpc error",
			err:  status.Error(codes.DataLoss, "dummy"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var wg sync.WaitGroup
			wg.Add(1)
			var md metadata.MD
			ctx := metadata.NewIncomingContext(context.Background(), md)
			stream := mockedServerStream{ctx: ctx}
			err := StreamTracingInterceptor(nil, &stream, &grpc.StreamServerInfo{
				FullMethod: "/foo",
			}, func(svr interface{}, stream grpc.ServerStream) error {
				defer wg.Done()
				return test.err
			})
			wg.Wait()
			assert.Equal(t, test.err, err)
		})
	}
}

func TestClientStream_RecvMsg(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "nil error",
		},
		{
			name: "EOF",
			err:  io.EOF,
		},
		{
			name: "dummy error",
			err:  errors.New("dummy"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			stream := wrapServerStream(context.Background(), &mockedServerStream{
				ctx: context.Background(),
				err: test.err,
			})
			assert.Equal(t, test.err, stream.RecvMsg(nil))
		})
	}
}

func TestServerStream_SendMsg(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "nil error",
		},
		{
			name: "with error",
			err:  errors.New("dummy"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			stream := wrapServerStream(context.Background(), &mockedServerStream{
				ctx: context.Background(),
				err: test.err,
			})
			assert.Equal(t, test.err, stream.SendMsg(nil))
		})
	}
}

type mockedServerStream struct {
	ctx context.Context
	err error
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
	return m.err
}

func (m *mockedServerStream) RecvMsg(v interface{}) error {
	return m.err
}
