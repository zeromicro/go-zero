package serverinterceptors

import (
	"context"

	"github.com/zeromicro/go-zero/core/md"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryMdInterceptor returns an interceptor that can get the md.Metadata.
func UnaryMdInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx = extractMd(ctx)
	return handler(ctx, req)
}

// StreamMdInterceptor returns an interceptor that can get the md.Metadata.
func StreamMdInterceptor(svr interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	ctx := extractMd(ss.Context())
	return handler(svr, &wrappedServerStream{ServerStream: ss, ctx: ctx})
}

func extractMd(ctx context.Context) context.Context {
	incomingMd, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	ctx = md.Extract(ctx, md.GRPCMetadataCarrier(incomingMd))

	return ctx
}

var _ grpc.ServerStream = (*wrappedServerStream)(nil)

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
