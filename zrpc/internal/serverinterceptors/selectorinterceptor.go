package serverinterceptors

import (
	"context"

	"github.com/zeromicro/go-zero/core/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnarySelectorInterceptor returns an interceptor that can get the selector and colors.
func UnarySelectorInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx = extractMd(ctx)

	return handler(ctx, req)
}

// StreamSelectorInterceptor returns an interceptor that can get the selector and colors.
func StreamSelectorInterceptor(svr interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	ctx := extractMd(ss.Context())

	return handler(svr, &wrappedServerStream{ss: ss, ctx: ctx})
}

func extractMd(ctx context.Context) context.Context {
	incomingMd, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	selectorVal := incomingMd.Get("selector")
	if len(selectorVal) != 0 {
		ctx = selector.NewSelectorContext(ctx, selectorVal[0])
	}

	colorsVal := incomingMd.Get("colors")
	if len(colorsVal) != 0 {
		ctx = selector.NewColorsContext(ctx, colorsVal...)
	}

	return ctx
}

var _ grpc.ServerStream = (*wrappedServerStream)(nil)

type wrappedServerStream struct {
	ss  grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) SetHeader(m metadata.MD) error {
	return w.ss.SetHeader(m)
}

func (w *wrappedServerStream) SendHeader(m metadata.MD) error {
	return w.ss.SendHeader(m)
}

func (w *wrappedServerStream) SetTrailer(m metadata.MD) {
	w.ss.SetTrailer(m)
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func (w *wrappedServerStream) SendMsg(m interface{}) error {
	return w.ss.SendMsg(m)
}

func (w *wrappedServerStream) RecvMsg(m interface{}) error {
	return w.ss.RecvMsg(m)
}
