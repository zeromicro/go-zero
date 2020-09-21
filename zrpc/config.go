package zrpc

import (
	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/core/service"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

type (
	RpcServerConf struct {
		service.ServiceConf
		ListenOn      string
		Etcd          discov.EtcdConf    `json:",optional"`
		Auth          bool               `json:",optional"`
		Redis         redis.RedisKeyConf `json:",optional"`
		StrictControl bool               `json:",optional"`
		// pending forever is not allowed
		// never set it to 0, if zero, the underlying will set to 2s automatically
		Timeout      int64 `json:",default=2000"`
		CpuThreshold int64 `json:",default=900,range=[0:1000]"`
	}

	RpcClientConf struct {
		Etcd      discov.EtcdConf `json:",optional"`
		Endpoints []string        `json:",optional=!Etcd"`
		App       string          `json:",optional"`
		Token     string          `json:",optional"`
		Timeout   int64           `json:",optional"`
	}
)

func NewDirectClientConf(endpoints []string, app, token string) RpcClientConf {
	return RpcClientConf{
		Endpoints: endpoints,
		App:       app,
		Token:     token,
	}
}

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

func (sc RpcServerConf) HasEtcd() bool {
	return len(sc.Etcd.Hosts) > 0 && len(sc.Etcd.Key) > 0
}

func (sc RpcServerConf) Validate() error {
	if sc.Auth {
		if err := sc.Redis.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (cc RpcClientConf) HasCredential() bool {
	return len(cc.App) > 0 && len(cc.Token) > 0
}
