package internal

import (
	"github.com/zeromicro/go-zero/core/stat"
	"google.golang.org/grpc"
)

type (
	// RegisterFn defines the method to register a server.
	RegisterFn func(*grpc.Server)

	// Server interface represents a rpc server.
	Server interface {
		AddOptions(options ...grpc.ServerOption)
		AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor)
		AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor)
		SetName(string)
		Start(register RegisterFn) error
	}

	baseRpcServer struct {
		address            string
		metrics            *stat.Metrics
		options            []grpc.ServerOption
		streamInterceptors []grpc.StreamServerInterceptor
		unaryInterceptors  []grpc.UnaryServerInterceptor
	}
)

func newBaseRpcServer(address string, rpcServerOpts *rpcServerOptions) *baseRpcServer {
	return &baseRpcServer{
		address: address,
		metrics: rpcServerOpts.metrics,
	}
}

func (s *baseRpcServer) AddOptions(options ...grpc.ServerOption) {
	s.options = append(s.options, options...)
}

func (s *baseRpcServer) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	s.streamInterceptors = append(s.streamInterceptors, interceptors...)
}

func (s *baseRpcServer) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	s.unaryInterceptors = append(s.unaryInterceptors, interceptors...)
}

func (s *baseRpcServer) SetName(name string) {
	s.metrics.SetName(name)
}
