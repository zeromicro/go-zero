package zrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestRpcClientConf(t *testing.T) {
	t.Run("direct", func(t *testing.T) {
		conf := NewDirectClientConf([]string{"localhost:1234"}, "foo", "bar")
		assert.True(t, conf.HasCredential())
	})

	t.Run("etcd", func(t *testing.T) {
		conf := NewEtcdClientConf([]string{"localhost:1234", "localhost:5678"},
			"key", "foo", "bar")
		assert.True(t, conf.HasCredential())
	})

	t.Run("etcd with account", func(t *testing.T) {
		conf := NewEtcdClientConf([]string{"localhost:1234", "localhost:5678"},
			"key", "foo", "bar")
		conf.Etcd.User = "user"
		conf.Etcd.Pass = "pass"
		_, err := conf.BuildTarget()
		assert.NoError(t, err)
	})

	t.Run("etcd with tls", func(t *testing.T) {
		conf := NewEtcdClientConf([]string{"localhost:1234", "localhost:5678"},
			"key", "foo", "bar")
		conf.Etcd.CertFile = "cert"
		conf.Etcd.CertKeyFile = "key"
		conf.Etcd.CACertFile = "ca"
		_, err := conf.BuildTarget()
		assert.Error(t, err)
	})
}

func TestRpcServerConf(t *testing.T) {
	conf := RpcServerConf{
		ServiceConf: service.ServiceConf{},
		ListenOn:    "",
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost:1234"},
			Key:   "key",
		},
		Auth: true,
		Redis: redis.RedisKeyConf{
			RedisConf: redis.RedisConf{
				Type: redis.NodeType,
			},
			Key: "foo",
		},
		StrictControl: false,
		Timeout:       0,
		CpuThreshold:  0,
	}
	assert.True(t, conf.HasEtcd())
	assert.NotNil(t, conf.Validate())
	conf.Redis.Host = "localhost:5678"
	assert.Nil(t, conf.Validate())
}
