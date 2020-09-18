package internal

import (
	"context"

	"google.golang.org/grpc"
)

func WithStreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.StreamInterceptor(chainStreamServerInterceptors(interceptors...))
}

func WithUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.UnaryInterceptor(chainUnaryServerInterceptors(interceptors...))
}

func chainStreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	switch len(interceptors) {
	case 0:
		return func(srv interface{}, stream grpc.ServerStream, _ *grpc.StreamServerInfo,
			handler grpc.StreamHandler) error {
			return handler(srv, stream)
		}
	case 1:
		return interceptors[0]
	default:
		last := len(interceptors) - 1
		return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo,
			handler grpc.StreamHandler) error {
			var chainHandler grpc.StreamHandler
			var current int

			chainHandler = func(curSrv interface{}, curStream grpc.ServerStream) error {
				if current == last {
					return handler(curSrv, curStream)
				}

				current++
				err := interceptors[current](curSrv, curStream, info, chainHandler)
				current--

				return err
			}

			return interceptors[0](srv, stream, info, chainHandler)
		}
	}
}

func chainUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	switch len(interceptors) {
	case 0:
		return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
			interface{}, error) {
			return handler(ctx, req)
		}
	case 1:
		return interceptors[0]
	default:
		last := len(interceptors) - 1
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
			interface{}, error) {
			var chainHandler grpc.UnaryHandler
			var current int

			chainHandler = func(curCtx context.Context, curReq interface{}) (interface{}, error) {
				if current == last {
					return handler(curCtx, curReq)
				}

				current++
				resp, err := interceptors[current](curCtx, curReq, info, chainHandler)
				current--

				return resp, err
			}

			return interceptors[0](ctx, req, info, chainHandler)
		}
	}
}
