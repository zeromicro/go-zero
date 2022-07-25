package clientinterceptors

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/md"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryMdInterceptor returns an interceptor that can inject md.Metadata.
func UnaryMdInterceptor(metadata md.Metadata) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = injectionMd(ctx, metadata)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamMdInterceptor returns an interceptor that can inject md.Metadata.
func StreamMdInterceptor(metadata md.Metadata) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = injectionMd(ctx, metadata)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func injectionMd(ctx context.Context, defaultMetadata md.Metadata) context.Context {
	m := md.FromContext(ctx)
	m = m.Clone()
	for key, values := range defaultMetadata {
		m.Append(key, values...)
	}

	ctx = md.NewContext(ctx, m)
	mdBytes, err := json.Marshal(m)
	if err != nil {
		logx.WithContext(ctx).Error(err)
	} else {
		ctx = appendToOutgoingContext(ctx, "metadata", string(mdBytes))
	}
	grpcMd := metadata.MD{}
	md.Inject(ctx, md.GRPCMetadataCarrier(grpcMd))

	return ctx
}

func appendToOutgoingContext(ctx context.Context, key string, value string) context.Context {
	if value == "" {
		return ctx
	}

	m, b := metadata.FromOutgoingContext(ctx)
	if !b {
		m = metadata.MD{}
	} else {
		m = m.Copy()
	}

	m.Append(key, value)

	return metadata.NewOutgoingContext(ctx, m)
}
