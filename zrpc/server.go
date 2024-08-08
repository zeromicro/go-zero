package zrpc

import (
	"time"

	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc/internal"
	"github.com/zeromicro/go-zero/zrpc/internal/auth"
	"github.com/zeromicro/go-zero/zrpc/internal/serverinterceptors"
	"google.golang.org/grpc"
)

// A RpcServer is a rpc server.
type RpcServer struct {
	server   internal.Server
	register internal.RegisterFn
}

// MustNewServer returns a RpcSever, exits on any error.
func MustNewServer(c RpcServerConf, register internal.RegisterFn) *RpcServer {
	server, err := NewServer(c, register)
	logx.Must(err)
	return server
}

// NewServer returns a RpcServer.
func NewServer(c RpcServerConf, register internal.RegisterFn) (*RpcServer, error) {
	var err error
	if err = c.Validate(); err != nil {
		return nil, err
	}

	var server internal.Server
	metrics := stat.NewMetrics(c.ListenOn)
	serverOptions := []internal.ServerOption{
		internal.WithRpcHealth(c.Health),
	}

	if c.HasEtcd() {
		server, err = internal.NewRpcPubServer(c.Etcd, c.ListenOn, serverOptions...)
		if err != nil {
			return nil, err
		}
	} else {
		server = internal.NewRpcServer(c.ListenOn, serverOptions...)
	}

	server.SetName(c.Name)
	metrics.SetName(c.Name)
	setupStreamInterceptors(server, c)
	setupUnaryInterceptors(server, c, metrics)
	if err = setupAuthInterceptors(server, c); err != nil {
		return nil, err
	}

	rpcServer := &RpcServer{
		server:   server,
		register: register,
	}
	if err = c.SetUp(); err != nil {
		return nil, err
	}

	return rpcServer, nil
}

// AddOptions adds given options.
func (rs *RpcServer) AddOptions(options ...grpc.ServerOption) {
	rs.server.AddOptions(options...)
}

// AddStreamInterceptors adds given stream interceptors.
func (rs *RpcServer) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	rs.server.AddStreamInterceptors(interceptors...)
}

// AddUnaryInterceptors adds given unary interceptors.
func (rs *RpcServer) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	rs.server.AddUnaryInterceptors(interceptors...)
}

// Start starts the RpcServer.
// Graceful shutdown is enabled by default.
// Use proc.SetTimeToForceQuit to customize the graceful shutdown period.
func (rs *RpcServer) Start() {
	if err := rs.server.Start(rs.register); err != nil {
		logx.Error(err)
		panic(err)
	}
}

// Stop stops the RpcServer.
func (rs *RpcServer) Stop() {
	logx.Close()
}

// DontLogContentForMethod disable logging content for given method.
// Deprecated: use ServerMiddlewaresConf.IgnoreContentMethods instead.
func DontLogContentForMethod(method string) {
	serverinterceptors.DontLogContentForMethod(method)
}

// SetServerSlowThreshold sets the slow threshold on server side.
// Deprecated: use ServerMiddlewaresConf.SlowThreshold instead.
func SetServerSlowThreshold(threshold time.Duration) {
	serverinterceptors.SetSlowThreshold(threshold)
}

func setupAuthInterceptors(svr internal.Server, c RpcServerConf) error {
	if !c.Auth {
		return nil
	}
	rds, err := redis.NewRedis(c.Redis.RedisConf)
	if err != nil {
		return err
	}

	authenticator, err := auth.NewAuthenticator(rds, c.Redis.Key, c.StrictControl)
	if err != nil {
		return err
	}

	svr.AddStreamInterceptors(serverinterceptors.StreamAuthorizeInterceptor(authenticator))
	svr.AddUnaryInterceptors(serverinterceptors.UnaryAuthorizeInterceptor(authenticator))

	return nil
}

func setupStreamInterceptors(svr internal.Server, c RpcServerConf) {
	if c.Middlewares.Trace {
		svr.AddStreamInterceptors(serverinterceptors.StreamTracingInterceptor)
	}
	if c.Middlewares.Recover {
		svr.AddStreamInterceptors(serverinterceptors.StreamRecoverInterceptor)
	}
	if c.Middlewares.Breaker {
		svr.AddStreamInterceptors(serverinterceptors.StreamBreakerInterceptor)
	}
}

func setupUnaryInterceptors(svr internal.Server, c RpcServerConf, metrics *stat.Metrics) {
	if c.Middlewares.Trace {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryTracingInterceptor)
	}
	if c.Middlewares.Recover {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryRecoverInterceptor)
	}
	if c.Middlewares.Stat {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryStatInterceptor(metrics, c.Middlewares.StatConf))
	}
	if c.Middlewares.Prometheus {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryPrometheusInterceptor)
	}
	if c.Middlewares.Breaker {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryBreakerInterceptor)
	}
	if c.CpuThreshold > 0 {
		shedder := load.NewAdaptiveShedder(load.WithCpuThreshold(c.CpuThreshold))
		svr.AddUnaryInterceptors(serverinterceptors.UnarySheddingInterceptor(shedder, metrics))
	}
	if c.Timeout > 0 {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryTimeoutInterceptor(
			time.Duration(c.Timeout)*time.Millisecond, c.MethodTimeouts...))
	}
}
