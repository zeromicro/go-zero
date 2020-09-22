package internal

import (
	"context"

	"google.golang.org/grpc"
)

func WithStreamClientInterceptors(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithStreamInterceptor(chainStreamClientInterceptors(interceptors...))
}

func WithUnaryClientInterceptors(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithUnaryInterceptor(chainUnaryClientInterceptors(interceptors...))
}

func chainStreamClientInterceptors(interceptors ...grpc.StreamClientInterceptor) grpc.StreamClientInterceptor {
	switch len(interceptors) {
	case 0:
		return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
			streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			return streamer(ctx, desc, cc, method, opts...)
		}
	case 1:
		return interceptors[0]
	default:
		last := len(interceptors) - 1
		return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
			method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			var chainStreamer grpc.Streamer
			var current int

			chainStreamer = func(curCtx context.Context, curDesc *grpc.StreamDesc, curCc *grpc.ClientConn,
				curMethod string, curOpts ...grpc.CallOption) (grpc.ClientStream, error) {
				if current == last {
					return streamer(curCtx, curDesc, curCc, curMethod, curOpts...)
				}

				current++
				clientStream, err := interceptors[current](curCtx, curDesc, curCc, curMethod, chainStreamer, curOpts...)
				current--

				return clientStream, err
			}

			return interceptors[0](ctx, desc, cc, method, chainStreamer, opts...)
		}
	}
}

func chainUnaryClientInterceptors(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	switch len(interceptors) {
	case 0:
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	case 1:
		return interceptors[0]
	default:
		last := len(interceptors) - 1
		return func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			var chainInvoker grpc.UnaryInvoker
			var current int

			chainInvoker = func(curCtx context.Context, curMethod string, curReq, curReply interface{},
				curCc *grpc.ClientConn, curOpts ...grpc.CallOption) error {
				if current == last {
					return invoker(curCtx, curMethod, curReq, curReply, curCc, curOpts...)
				}

				current++
				err := interceptors[current](curCtx, curMethod, curReq, curReply, curCc, chainInvoker, curOpts...)
				current--

				return err
			}

			return interceptors[0](ctx, method, req, reply, cc, chainInvoker, opts...)
		}
	}
}
