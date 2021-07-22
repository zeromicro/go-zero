package internal

import "google.golang.org/grpc"

// WithStreamClientInterceptors uses given client stream interceptors.
func WithStreamClientInterceptors(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(interceptors...)
}

// WithUnaryClientInterceptors uses given client unary interceptors.
func WithUnaryClientInterceptors(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}
