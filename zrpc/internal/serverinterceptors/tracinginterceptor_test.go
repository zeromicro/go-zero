package serverinterceptors

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	ztrace "github.com/zeromicro/go-zero/core/trace"
	"github.com/zeromicro/go-zero/core/trace/tracetest"
	"go.opentelemetry.io/otel/attribute"
	tcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestUnaryOpenTracingInterceptor_Disable(t *testing.T) {
	_, err := UnaryTracingInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req any) (any, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestUnaryOpenTracingInterceptor_Enabled(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})
	defer ztrace.StopAgent()

	_, err := UnaryTracingInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/package.TestService.GetUser",
	}, func(ctx context.Context, req any) (any, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestUnaryTracingInterceptor(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		var run int32
		me := tracetest.NewInMemoryExporter(t)
		_, err := UnaryTracingInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/proto.Hello/Echo",
		}, func(ctx context.Context, req any) (any, error) {
			atomic.AddInt32(&run, 1)
			return nil, nil
		})
		assert.Nil(t, err)
		assert.Equal(t, int32(1), atomic.LoadInt32(&run))

		assert.Equal(t, 1, len(me.GetSpans()))
		span := me.GetSpans()[0].Snapshot()
		assert.Equal(t, 2, len(span.Events()))
		assert.ElementsMatch(t, []attribute.KeyValue{
			ztrace.RPCSystemGRPC,
			semconv.RPCServiceKey.String("proto.Hello"),
			semconv.RPCMethodKey.String("Echo"),
			ztrace.StatusCodeAttr(codes.OK),
		}, span.Attributes())
	})

	t.Run("grpc error status", func(t *testing.T) {
		me := tracetest.NewInMemoryExporter(t)
		_, err := UnaryTracingInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/proto.Hello/Echo",
		}, func(ctx context.Context, req any) (any, error) {
			return nil, status.Errorf(codes.Unknown, "test")
		})
		assert.Error(t, err)
		assert.Equal(t, 1, len(me.GetSpans()))
		span := me.GetSpans()[0].Snapshot()
		assert.Equal(t, trace.Status{
			Code:        tcodes.Error,
			Description: "test",
		}, span.Status())
		assert.Equal(t, 2, len(span.Events()))
		assert.ElementsMatch(t, []attribute.KeyValue{
			ztrace.RPCSystemGRPC,
			semconv.RPCServiceKey.String("proto.Hello"),
			semconv.RPCMethodKey.String("Echo"),
			ztrace.StatusCodeAttr(codes.Unknown),
		}, span.Attributes())
	})

	t.Run("non grpc status error", func(t *testing.T) {
		me := tracetest.NewInMemoryExporter(t)
		_, err := UnaryTracingInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/proto.Hello/Echo",
		}, func(ctx context.Context, req any) (any, error) {
			return nil, errors.New("test")
		})
		assert.Error(t, err)
		assert.Equal(t, 1, len(me.GetSpans()))
		span := me.GetSpans()[0].Snapshot()
		assert.Equal(t, trace.Status{
			Code:        tcodes.Error,
			Description: "test",
		}, span.Status())
		assert.Equal(t, 1, len(span.Events()))
		assert.ElementsMatch(t, []attribute.KeyValue{
			ztrace.RPCSystemGRPC,
			semconv.RPCServiceKey.String("proto.Hello"),
			semconv.RPCMethodKey.String("Echo"),
		}, span.Attributes())
	})
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
			}, func(ctx context.Context, req any) (any, error) {
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
	}, func(svr any, stream grpc.ServerStream) error {
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
			}, func(svr any, stream grpc.ServerStream) error {
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
			}, func(svr any, stream grpc.ServerStream) error {
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

func (m *mockedServerStream) SetHeader(_ metadata.MD) error {
	panic("implement me")
}

func (m *mockedServerStream) SendHeader(_ metadata.MD) error {
	panic("implement me")
}

func (m *mockedServerStream) SetTrailer(_ metadata.MD) {
	panic("implement me")
}

func (m *mockedServerStream) Context() context.Context {
	if m.ctx == nil {
		return context.Background()
	}

	return m.ctx
}

func (m *mockedServerStream) SendMsg(_ any) error {
	return m.err
}

func (m *mockedServerStream) RecvMsg(_ any) error {
	return m.err
}
