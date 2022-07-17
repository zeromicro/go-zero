package clientinterceptors

import (
	"context"

	"github.com/zeromicro/go-zero/core/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnarySelectorInterceptor returns an interceptor that can inject selector and colors.
func UnarySelectorInterceptor(defaultSelectorName string, defaultColors []string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = injectionSelectorName(ctx, defaultSelectorName, defaultColors)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamSelectorInterceptor returns an interceptor that can inject selector and colors.
func StreamSelectorInterceptor(defaultSelectorName string, defaultColors []string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = injectionSelectorName(ctx, defaultSelectorName, defaultColors)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func injectionSelectorName(ctx context.Context, defaultSelectorName string, defaultColors []string) context.Context {
	selectorName, ok := selector.SelectorFromContext(ctx)
	if !ok && defaultSelectorName != "" {
		selectorName = defaultSelectorName
		ctx = selector.NewSelectorContext(ctx, selectorName)
	}
	if selectorName != "" {
		ctx = appendToOutgoingContext(ctx, "selector", selectorName)
	}

	var colors []string
	c, ok := selector.ColorsFromContext(ctx)
	if !ok {
		if len(defaultColors) != 0 {
			colors = defaultColors
			ctx = selector.NewColorsContext(ctx, colors...)
		}
	} else {
		colors = c.Colors()
	}
	if len(colors) != 0 {
		ctx = appendToOutgoingContext(ctx, "colors", colors...)
	}

	return ctx
}

func appendToOutgoingContext(ctx context.Context, key string, values ...string) context.Context {
	md, b := metadata.FromOutgoingContext(ctx)
	if !b {
		md = metadata.MD{}
	} else {
		md = md.Copy()
	}

	md.Set(key, values...)
	return metadata.NewOutgoingContext(ctx, md)
}
