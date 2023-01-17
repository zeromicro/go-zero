package internal

import (
	"fmt"
	"net"

	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/internal/health"
	"github.com/zeromicro/go-zero/zrpc/internal/serverinterceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const probeNamePrefix = "zrpc"

type (
	// ServerOption defines the method to customize a rpcServerOptions.
	ServerOption func(options *rpcServerOptions)

	rpcServerOptions struct {
		metrics *stat.Metrics
		health  bool
	}

	rpcServer struct {
		*baseRpcServer
		name          string
		middlewares   ServerMiddlewaresConf
		healthManager health.Probe
	}
)

// NewRpcServer returns a Server.
func NewRpcServer(addr string, middlewares ServerMiddlewaresConf, opts ...ServerOption) Server {
	var options rpcServerOptions
	for _, opt := range opts {
		opt(&options)
	}
	if options.metrics == nil {
		options.metrics = stat.NewMetrics(addr)
	}

	return &rpcServer{
		baseRpcServer: newBaseRpcServer(addr, &options),
		middlewares:   middlewares,
		healthManager: health.NewHealthManager(fmt.Sprintf("%s-%s", probeNamePrefix, addr)),
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

	unaryInterceptors := s.buildUnaryInterceptors()
	unaryInterceptors = append(unaryInterceptors, s.unaryInterceptors...)
	streamInterceptors := s.buildStreamInterceptors()
	streamInterceptors = append(streamInterceptors, s.streamInterceptors...)
	options := append(s.options, grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...))
	server := grpc.NewServer(options...)
	register(server)

	// register the health check service
	if s.health != nil {
		grpc_health_v1.RegisterHealthServer(server, s.health)
		s.health.Resume()
	}
	s.healthManager.MarkReady()
	health.AddProbe(s.healthManager)

	// we need to make sure all others are wrapped up,
	// so we do graceful stop at shutdown phase instead of wrap up phase
	waitForCalled := proc.AddWrapUpListener(func() {
		if s.health != nil {
			s.health.Shutdown()
		}
		server.GracefulStop()
	})
	defer waitForCalled()

	return server.Serve(lis)
}

func (s *rpcServer) buildStreamInterceptors() []grpc.StreamServerInterceptor {
	var interceptors []grpc.StreamServerInterceptor

	if s.middlewares.Trace {
		interceptors = append(interceptors, serverinterceptors.StreamTracingInterceptor)
	}
	if s.middlewares.Recover {
		interceptors = append(interceptors, serverinterceptors.StreamRecoverInterceptor)
	}
	if s.middlewares.Breaker {
		interceptors = append(interceptors, serverinterceptors.StreamBreakerInterceptor)
	}

	return interceptors
}

func (s *rpcServer) buildUnaryInterceptors() []grpc.UnaryServerInterceptor {
	var interceptors []grpc.UnaryServerInterceptor

	if s.middlewares.Trace {
		interceptors = append(interceptors, serverinterceptors.UnaryTracingInterceptor)
	}
	if s.middlewares.Recover {
		interceptors = append(interceptors, serverinterceptors.UnaryRecoverInterceptor)
	}
	if s.middlewares.Stat {
		interceptors = append(interceptors, serverinterceptors.UnaryStatInterceptor(s.metrics))
	}
	if s.middlewares.Prometheus {
		interceptors = append(interceptors, serverinterceptors.UnaryPrometheusInterceptor)
	}
	if s.middlewares.Breaker {
		interceptors = append(interceptors, serverinterceptors.UnaryBreakerInterceptor)
	}

	return interceptors
}

// WithMetrics returns a func that sets metrics to a Server.
func WithMetrics(metrics *stat.Metrics) ServerOption {
	return func(options *rpcServerOptions) {
		options.metrics = metrics
	}
}

// WithRpcHealth returns a func that sets rpc health switch to a Server.
func WithRpcHealth(health bool) ServerOption {
	return func(options *rpcServerOptions) {
		options.health = health
	}
}
