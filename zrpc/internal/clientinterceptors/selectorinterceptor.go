package clientinterceptors

import (
	"context"

	"github.com/zeromicro/go-zero/core/selector"
	"google.golang.org/grpc"
)

// UnarySelectorInterceptor returns an interceptor that can inject selector.
func UnarySelectorInterceptor(defaultSelectorName string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = injectionSelector(ctx, defaultSelectorName)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamSelectorInterceptor returns an interceptor that can inject selector.
func StreamSelectorInterceptor(defaultSelectorName string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = injectionSelector(ctx, defaultSelectorName)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func injectionSelector(ctx context.Context, defaultSelectorName string) context.Context {
	selectorName := selector.FromContext(ctx)
	if selectorName == "" {
		selectorName = defaultSelectorName
	}
	ctx = selector.NewContext(ctx, selectorName)
	ctx = appendToOutgoingContext(ctx, "selector", selectorName)

	return ctx
}
