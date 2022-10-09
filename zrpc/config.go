package zrpc

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc/resolver"
)

type (
	// A RpcServerConf is a rpc server config.
	RpcServerConf struct {
		service.ServiceConf `yaml:",inline"`
		ListenOn            string             `json:"ListenOn" yaml:"ListenOn"`
		Etcd                discov.EtcdConf    `json:"Etcd,optional" yaml:"Etcd"`
		Auth                bool               `json:"Auth,optional" yaml:"Auth"`
		Redis               redis.RedisKeyConf `json:"Redis,optional" yaml:"Redis"`
		StrictControl       bool               `json:"StrictControl,optional" yaml:"StrictControl"`
		Timeout             int64              `json:"Timeout,default=2000" yaml:"Timeout"` // setting 0 means no timeout
		CpuThreshold        int64              `json:"CpuThreshold,default=900,range=[0:1000]" yaml:"CpuThreshold"`
		Health              bool               `json:"Health,default=true" yaml:"Health"` // grpc health check switch
	}

	// A RpcClientConf is a rpc client config.
	RpcClientConf struct {
		Etcd      discov.EtcdConf `json:"Etcd,optional" yaml:"Etcd"`
		Endpoints []string        `json:"Endpoints,optional" yaml:"Endpoints"`
		Target    string          `json:"Target,optional" yaml:"Target"`
		App       string          `json:"App,optional" yaml:"App"`
		Token     string          `json:"Token,optional" yaml:"Token"`
		NonBlock  bool            `json:"NonBlock,optional" yaml:"NonBlock"`
		Timeout   int64           `json:"Timeout,default=2000" yaml:"Timeout"`
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
