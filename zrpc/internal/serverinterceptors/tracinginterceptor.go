package serverinterceptors

import (
	"context"

	"github.com/tal-tech/go-zero/core/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryTracingInterceptor returns a grpc.UnaryServerInterceptor
// that handles tracing with given service name.
func UnaryTracingInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		carrier, err := trace.Extract(trace.GrpcFormat, md)
		if err != nil {
			return handler(ctx, req)
		}

		ctx, span := trace.StartServerSpan(ctx, carrier, serviceName, info.FullMethod)
		defer span.Finish()
		return handler(ctx, req)
	}
}

// StreamTracingInterceptor returns a grpc.StreamServerInterceptor
// that handles tracing with given service name.
func StreamTracingInterceptor(serviceName string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo,
		handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(srv, ss)
		}

		carrier, err := trace.Extract(trace.GrpcFormat, md)
		if err != nil {
			return handler(srv, ss)
		}

		ctx, span := trace.StartServerSpan(ctx, carrier, serviceName, info.FullMethod)
		defer span.Finish()
		return handler(srv, ss)
	}
}
