package redis

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
)

func TestRedisConfTLSVerification(t *testing.T) {
	secure := RedisConf{
		Host: "redis.example.com:6379",
		Type: NodeType,
		Tls:  true,
	}.NewRedis()
	assert.False(t, secure.tlsConfig.InsecureSkipVerify)

	insecure := RedisConf{
		Host:                  "redis.example.com:6379",
		Type:                  NodeType,
		Tls:                   true,
		TlsInsecureSkipVerify: true,
	}.NewRedis()
	assert.True(t, insecure.tlsConfig.InsecureSkipVerify)
	insecureAgain := RedisConf{
		Host:                  "redis.example.com:6379",
		Type:                  NodeType,
		Tls:                   true,
		TlsInsecureSkipVerify: true,
	}.NewRedis()
	assert.Equal(t, insecure.tlsConfigKey, insecureAgain.tlsConfigKey)
	insecurePrimary, err := NewRedis(RedisConf{
		Host:                  "redis.example.com:6379",
		Type:                  NodeType,
		Tls:                   true,
		TlsInsecureSkipVerify: true,
		NonBlock:              true,
	})
	assert.NoError(t, err)
	assert.True(t, insecurePrimary.tlsConfig.InsecureSkipVerify)

	custom, err := NewRedis(RedisConf{
		Host:     "redis.example.com:6379",
		Type:     NodeType,
		Tls:      true,
		NonBlock: true,
	}, WithTLSConfig(&tls.Config{ServerName: "redis.internal"}))
	assert.NoError(t, err)
	assert.Equal(t, "redis.internal", custom.tlsConfig.ServerName)
}

func TestRedisConf(t *testing.T) {
	tests := []struct {
		name string
		RedisConf
		ok bool
	}{
		{
			name: "missing host",
			RedisConf: RedisConf{
				Host: "",
				Type: NodeType,
				Pass: "",
			},
			ok: false,
		},
		{
			name: "missing type",
			RedisConf: RedisConf{
				Host: "localhost:6379",
				Type: "",
				Pass: "",
			},
			ok: false,
		},
		{
			name: "ok",
			RedisConf: RedisConf{
				Host: "localhost:6379",
				Type: NodeType,
				Pass: "",
			},
			ok: true,
		},
		{
			name: "ok",
			RedisConf: RedisConf{
				Host: "localhost:6379",
				Type: ClusterType,
				Pass: "pwd",
				Tls:  true,
			},
			ok: true,
		},
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			if test.ok {
				assert.Nil(t, test.RedisConf.Validate())
				assert.NotNil(t, test.RedisConf.NewRedis())
			} else {
				assert.NotNil(t, test.RedisConf.Validate())
			}
		})
	}
}

func TestRedisKeyConf(t *testing.T) {
	tests := []struct {
		name string
		RedisKeyConf
		ok bool
	}{
		{
			name: "missing host",
			RedisKeyConf: RedisKeyConf{
				RedisConf: RedisConf{
					Host: "",
					Type: NodeType,
					Pass: "",
				},
				Key: "foo",
			},
			ok: false,
		},
		{
			name: "missing key",
			RedisKeyConf: RedisKeyConf{
				RedisConf: RedisConf{
					Host: "localhost:6379",
					Type: NodeType,
					Pass: "",
				},
				Key: "",
			},
			ok: false,
		},
		{
			name: "ok",
			RedisKeyConf: RedisKeyConf{
				RedisConf: RedisConf{
					Host: "localhost:6379",
					Type: NodeType,
					Pass: "",
				},
				Key: "foo",
			},
			ok: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.ok {
				assert.Nil(t, test.RedisKeyConf.Validate())
			} else {
				assert.NotNil(t, test.RedisKeyConf.Validate())
			}
		})
	}
}
