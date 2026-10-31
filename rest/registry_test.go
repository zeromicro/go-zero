package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov"
)

func TestRestConf_HasEtcd(t *testing.T) {
	tests := []struct {
		name     string
		conf     RestConf
		expected bool
	}{
		{
			name: "has etcd config",
			conf: RestConf{
				Etcd: discov.EtcdConf{
					Hosts: []string{"localhost:2379"},
					Key:   "test-service",
				},
			},
			expected: true,
		},
		{
			name: "missing etcd hosts",
			conf: RestConf{
				Etcd: discov.EtcdConf{
					Key: "test-service",
				},
			},
			expected: false,
		},
		{
			name: "missing etcd key",
			conf: RestConf{
				Etcd: discov.EtcdConf{
					Hosts: []string{"localhost:2379"},
				},
			},
			expected: false,
		},
		{
			name: "no etcd config",
			conf: RestConf{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.conf.HasEtcd())
		})
	}
}

func TestFigureOutListenAddr(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     int
		expected string
	}{
		{
			name:     "host and port specified",
			host:     "127.0.0.1",
			port:     8080,
			expected: "127.0.0.1:8080",
		},
		{
			name:     "all eths host",
			host:     "0.0.0.0",
			port:     8080,
			expected: "127.0.0.1:8080", // Will use internal IP
		},
		{
			name:     "empty host",
			host:     "",
			port:     8080,
			expected: "127.0.0.1:8080", // Will use internal IP
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := figureOutListenAddr(tt.host, tt.port)
			// For cases where internal IP is used, we just check that port is correct
			if tt.host == "0.0.0.0" || tt.host == "" {
				assert.Contains(t, result, ":8080")
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
