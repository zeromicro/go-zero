package rpcx

import (
	"log"
	"os"
	"strings"
	"time"

	"zero/core/load"
	"zero/core/logx"
	"zero/core/netx"
	"zero/core/rpc"
	"zero/core/rpc/serverinterceptors"
	"zero/core/stat"
	"zero/rpcx/auth"
	"zero/rpcx/interceptors"
)

const envPodIp = "POD_IP"

type RpcServer struct {
	server   rpc.Server
	register rpc.RegisterFn
}

func MustNewServer(c RpcServerConf, register rpc.RegisterFn) *RpcServer {
	server, err := NewServer(c, register)
	if err != nil {
		log.Fatal(err)
	}

	return server
}

func NewServer(c RpcServerConf, register rpc.RegisterFn) (*RpcServer, error) {
	var err error
	if err = c.Validate(); err != nil {
		return nil, err
	}

	var server rpc.Server
	metrics := stat.NewMetrics(c.ListenOn)
	if c.HasEtcd() {
		listenOn := figureOutListenOn(c.ListenOn)
		server, err = rpc.NewRpcPubServer(c.Etcd.Hosts, c.Etcd.Key, listenOn, rpc.WithMetrics(metrics))
		if err != nil {
			return nil, err
		}
	} else {
		server = rpc.NewRpcServer(c.ListenOn, rpc.WithMetrics(metrics))
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

func (rs *RpcServer) Start() {
	if err := rs.server.Start(rs.register); err != nil {
		logx.Error(err)
		panic(err)
	}
}

func (rs *RpcServer) Stop() {
	logx.Close()
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != "0.0.0.0" {
		return listenOn
	}

	ip := os.Getenv(envPodIp)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	} else {
		return strings.Join(append([]string{ip}, fields[1:]...), ":")
	}
}

func setupInterceptors(server rpc.Server, c RpcServerConf, metrics *stat.Metrics) error {
	if c.CpuThreshold > 0 {
		shedder := load.NewAdaptiveShedder(load.WithCpuThreshold(c.CpuThreshold))
		server.AddUnaryInterceptors(serverinterceptors.UnarySheddingInterceptor(shedder, metrics))
	}

	if c.Timeout > 0 {
		server.AddUnaryInterceptors(serverinterceptors.UnaryTimeoutInterceptor(
			time.Duration(c.Timeout) * time.Millisecond))
	}

	server.AddUnaryInterceptors(serverinterceptors.UnaryTracingInterceptor(c.Name))

	if c.Auth {
		authenticator, err := auth.NewAuthenticator(c.Redis.NewRedis(), c.Redis.Key, c.StrictControl)
		if err != nil {
			return err
		}

		server.AddStreamInterceptors(interceptors.StreamAuthorizeInterceptor(authenticator))
		server.AddUnaryInterceptors(interceptors.UnaryAuthorizeInterceptor(authenticator))
	}

	return nil
}
