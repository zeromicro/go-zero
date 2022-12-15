package serverinterceptors

import (
	"context"

	"github.com/zeromicro/go-zero/core/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnarySelectorInterceptor returns an interceptor that can get the selector.
func UnarySelectorInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx = extractSelector(ctx)
	return handler(ctx, req)
}

// StreamSelectorInterceptor returns an interceptor that can get the selector.
func StreamSelectorInterceptor(svr interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	ctx := extractSelector(ss.Context())
	return handler(svr, &wrappedServerStream{ServerStream: ss, ctx: ctx})
}

func extractSelector(ctx context.Context) context.Context {
	incomingMd, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	selectorVal := incomingMd.Get("selector")
	if len(selectorVal) != 0 {
		ctx = selector.NewContext(ctx, selectorVal[0])
	}

	return ctx
}
