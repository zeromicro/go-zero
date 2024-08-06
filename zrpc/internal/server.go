package internal

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/keepalive"
)

const defaultConnectionIdleDuration = time.Minute * 5

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
		health             *health.Server
		options            []grpc.ServerOption
		streamInterceptors []grpc.StreamServerInterceptor
		unaryInterceptors  []grpc.UnaryServerInterceptor
	}
)

func newBaseRpcServer(address string, rpcServerOpts *rpcServerOptions) *baseRpcServer {
	var h *health.Server
	if rpcServerOpts.health {
		h = health.NewServer()
	}
	return &baseRpcServer{
		address: address,
		health:  h,
		options: []grpc.ServerOption{grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: defaultConnectionIdleDuration,
		})},
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
