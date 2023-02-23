package redis

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestRedisCluterConf_Validate(t *testing.T) {
	type fields struct {
		Hosts []string
		Pass  string
		Tls   bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "",
			fields: fields{
				Hosts: nil,
				Pass:  "",
				Tls:   false,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, ErrEmptyHosts)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisCluterConf{
				Hosts: tt.fields.Hosts,
				Pass:  tt.fields.Pass,
				Tls:   tt.fields.Tls,
			}
			tt.wantErr(t, r.Validate(), fmt.Sprintf("Validate()"))
		})
	}
}
