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
		internal.WithMetrics(metrics),
		internal.WithRpcHealth(c.Health),
	}

	if c.HasEtcd() {
		server, err = internal.NewRpcPubServer(c.Etcd, c.ListenOn, c.Middlewares, serverOptions...)
		if err != nil {
			return nil, err
		}
	} else {
		server = internal.NewRpcServer(c.ListenOn, c.Middlewares, serverOptions...)
	}

	server.SetName(c.Name)
	if err = setupInterceptors(server, c, metrics); err != nil {
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

func setupInterceptors(svr internal.Server, c RpcServerConf, metrics *stat.Metrics) error {
	if c.CpuThreshold > 0 {
		shedder := load.NewAdaptiveShedder(load.WithCpuThreshold(c.CpuThreshold))
		svr.AddUnaryInterceptors(serverinterceptors.UnarySheddingInterceptor(shedder, metrics))
	}

	if c.Timeout > 0 {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryTimeoutInterceptor(
			time.Duration(c.Timeout)*time.Millisecond, c.MethodTimeouts...))
	}

	if c.Auth {
		if err := setupAuthInterceptors(svr, c); err != nil {
			return err
		}
	}

	return nil
}
