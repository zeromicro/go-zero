package serverinterceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/metainfo"
)

// UnaryCustomKeysInterceptor extract custom keys from request metadata and inject it into request context and logger fields.
func UnaryCustomKeysInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// try to extract custom keys from request metadata
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			ctx = metainfo.CustomKeysMapPropagator.Extract(ctx, metainfo.GrpcHeaderCarrier(md))
		}

		// try to inject custom keys to logger
		if info := metainfo.GetMapFromContext(ctx); len(info) > 0 {
			ctx = logx.ContextWithFields(ctx, logx.Field("custom_keys", info))
		}

		return handler(ctx, req)
	}
}
