package clientinterceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/metainfo"
)

// UnaryCustomKeysInterceptor auto append custom keys data to grpc client request metadata.
func UnaryCustomKeysInterceptor(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var md metadata.MD
	// try to append custom metadata to client request metadata
	if m, ok := metadata.FromOutgoingContext(ctx); ok {
		md = m.Copy()
	} else {
		md = metadata.MD{}
	}

	metainfo.CustomKeysMapPropagator.Inject(ctx, metainfo.GrpcHeaderCarrier(md))
	ctx = metadata.NewOutgoingContext(ctx, md)

	return invoker(ctx, method, req, reply, cc, opts...)
}
