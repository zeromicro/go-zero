package zrpc

import (
	"time"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc/internal"
	"github.com/zeromicro/go-zero/zrpc/resolver"
)

type (
	// ClientMiddlewaresConf defines whether to use client middlewares.
	ClientMiddlewaresConf = internal.ClientMiddlewaresConf
	// ServerMiddlewaresConf defines whether to use server middlewares.
	ServerMiddlewaresConf = internal.ServerMiddlewaresConf
	// StatConf defines the stat config.
	StatConf = internal.StatConf
	// MethodTimeoutConf defines specified timeout for gRPC method.
	MethodTimeoutConf = internal.MethodTimeoutConf

	// A RpcClientConf is a rpc client config.
	RpcClientConf struct {
		Etcd          discov.EtcdConf `json:",optional,inherit"`
		Endpoints     []string        `json:",optional"`
		Target        string          `json:",optional"`
		App           string          `json:",optional"`
		Token         string          `json:",optional"`
		NonBlock      bool            `json:",optional"`
		Timeout       int64           `json:",default=2000"`
		KeepaliveTime time.Duration   `json:",optional"`
		Middlewares   ClientMiddlewaresConf
	}

	// A RpcServerConf is a rpc server config.
	RpcServerConf struct {
		service.ServiceConf
		ListenOn      string
		Etcd          discov.EtcdConf    `json:",optional,inherit"`
		Auth          bool               `json:",optional"`
		Redis         redis.RedisKeyConf `json:",optional"`
		StrictControl bool               `json:",optional"`
		// setting 0 means no timeout
		Timeout      int64 `json:",default=2000"`
		CpuThreshold int64 `json:",default=900,range=[0:1000)"`
		// grpc health check switch
		Health      bool `json:",default=true"`
		Middlewares ServerMiddlewaresConf
		// setting specified timeout for gRPC method
		MethodTimeouts []MethodTimeoutConf `json:",optional"`
	}
)

// NewDirectClientConf returns a RpcClientConf.
func NewDirectClientConf(endpoints []string, app, token string) RpcClientConf {
	return RpcClientConf{
		Endpoints: endpoints,
		App:       app,
		Token:     token,
	}
}

// NewEtcdClientConf returns a RpcClientConf.
func NewEtcdClientConf(hosts []string, key, app, token string) RpcClientConf {
	return RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: hosts,
			Key:   key,
		},
		App:   app,
		Token: token,
	}
}

// HasEtcd checks if there is etcd settings in config.
func (sc RpcServerConf) HasEtcd() bool {
	return len(sc.Etcd.Hosts) > 0 && len(sc.Etcd.Key) > 0
}

// Validate validates the config.
func (sc RpcServerConf) Validate() error {
	if !sc.Auth {
		return nil
	}

	return sc.Redis.Validate()
}

// BuildTarget builds the rpc target from the given config.
func (cc RpcClientConf) BuildTarget() (string, error) {
	if len(cc.Endpoints) > 0 {
		return resolver.BuildDirectTarget(cc.Endpoints), nil
	} else if len(cc.Target) > 0 {
		return cc.Target, nil
	}

	if err := cc.Etcd.Validate(); err != nil {
		return "", err
	}

	if cc.Etcd.HasAccount() {
		discov.RegisterAccount(cc.Etcd.Hosts, cc.Etcd.User, cc.Etcd.Pass)
	}
	if cc.Etcd.HasTLS() {
		if err := discov.RegisterTLS(cc.Etcd.Hosts, cc.Etcd.CertFile, cc.Etcd.CertKeyFile,
			cc.Etcd.CACertFile, cc.Etcd.InsecureSkipVerify); err != nil {
			return "", err
		}
	}

	return resolver.BuildDiscovTarget(cc.Etcd.Hosts, cc.Etcd.Key), nil
}

// HasCredential checks if there is a credential in config.
func (cc RpcClientConf) HasCredential() bool {
	return len(cc.App) > 0 && len(cc.Token) > 0
}
