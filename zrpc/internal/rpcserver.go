package internal

import (
	"net"

	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/zrpc/internal/serverinterceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type (
	// ServerOption defines the method to customize a rpcServerOptions.
	ServerOption func(options *rpcServerOptions)

	rpcServerOptions struct {
		metrics *stat.Metrics
	}

	rpcServer struct {
		name string
		*baseRpcServer
	}
)

func init() {
	InitLogger()
}

// NewRpcServer returns a Server.
func NewRpcServer(address string, opts ...ServerOption) Server {
	var options rpcServerOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics == nil {
		options.metrics = stat.NewMetrics(address)
	}

	return &rpcServer{
		baseRpcServer: newBaseRpcServer(address, &options),
	}
}

func (s *rpcServer) SetName(name string) {
	s.name = name
	s.baseRpcServer.SetName(name)
}

func (s *rpcServer) Start(register RegisterFn) error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		serverinterceptors.UnaryTracingInterceptor,
		serverinterceptors.UnaryCrashInterceptor,
		serverinterceptors.UnaryStatInterceptor(s.metrics),
		serverinterceptors.UnaryPrometheusInterceptor,
		serverinterceptors.UnaryBreakerInterceptor,
	}
	unaryInterceptors = append(unaryInterceptors, s.unaryInterceptors...)
	streamInterceptors := []grpc.StreamServerInterceptor{
		serverinterceptors.StreamTracingInterceptor,
		serverinterceptors.StreamCrashInterceptor,
		serverinterceptors.StreamBreakerInterceptor,
	}
	streamInterceptors = append(streamInterceptors, s.streamInterceptors...)
	options := append(s.options, WithUnaryServerInterceptors(unaryInterceptors...),
		WithStreamServerInterceptors(streamInterceptors...))
	server := grpc.NewServer(options...)
	register(server)

	// register the health check service
	grpc_health_v1.RegisterHealthServer(server, s.health)
	s.health.Resume()

	// we need to make sure all others are wrapped up,
	// so we do graceful stop at shutdown phase instead of wrap up phase
	waitForCalled := proc.AddWrapUpListener(func() {
		s.health.Shutdown()
		server.GracefulStop()
	})
	defer waitForCalled()

	return server.Serve(lis)
}

// WithMetrics returns a func that sets metrics to a Server.
func WithMetrics(metrics *stat.Metrics) ServerOption {
	return func(options *rpcServerOptions) {
		options.metrics = metrics
	}
}
