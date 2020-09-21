package internal

import (
	"github.com/tal-tech/go-zero/core/stat"
	"google.golang.org/grpc"
)

type (
	RegisterFn func(*grpc.Server)

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

func newBaseRpcServer(address string, metrics *stat.Metrics) *baseRpcServer {
	return &baseRpcServer{
		address: address,
		metrics: metrics,
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
