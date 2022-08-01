package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stringx"
)

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

func TestRedisConfDB(t *testing.T) {
	tests := []struct {
		name     string
		init     func() (RedisConf, error)
		expected uint
	}{
		{
			name: "db_unset",
			init: func() (RedisConf, error) {
				return RedisConf{
					Host: "localhost:6379",
					Type: NodeType,
				}, nil
			},
			expected: defaultDatabase,
		},
		{
			name: "db_unset2",
			init: func() (RedisConf, error) {
				var redisConf RedisConf
				err := conf.LoadFromJsonBytes([]byte(`{"Host":"localhost:6379","Type":"node"}`), &redisConf)
				if err != nil {
					return RedisConf{}, err
				}
				return redisConf, nil
			},
			expected: defaultDatabase,
		},
		{
			name: "db_set",
			init: func() (RedisConf, error) {
				return RedisConf{
					Host: "localhost:6379",
					Type: NodeType,
					DB:   1,
				}, nil
			},
			expected: 1,
		},
		{
			name: "db_set2",
			init: func() (RedisConf, error) {
				var redisConf RedisConf
				err := conf.LoadFromJsonBytes([]byte(`{"Host":"localhost:6379","Type":"node","DB":1}`), &redisConf)
				if err != nil {
					return RedisConf{}, err
				}
				return redisConf, nil
			},
			expected: 1,
		},
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			redisConfig, err := test.init()
			assert.NoError(t, err)
			assert.Equal(t, test.expected, redisConfig.DB)
		})
	}
}
