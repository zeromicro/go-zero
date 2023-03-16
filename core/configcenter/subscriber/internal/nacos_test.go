package internal

import (
	"testing"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/stretchr/testify/assert"
)

func TestNacosConf_Validate(t *testing.T) {
	assert.Error(t, NacosConf{}.Validate())
	assert.NoError(t, NacosConf{ServerConfigs: []ServerConfig{
		{
			Address: "127.0.0.1:8848",
		},
	}}.Validate())
}

func TestNacosConf_clientParam(t *testing.T) {
	tests := []struct {
		name    string
		conf    NacosConf
		wantErr bool
	}{
		{
			name:    "Validate fail",
			wantErr: true,
		},
		{
			name: "invalid address",
			conf: NacosConf{
				ServerConfigs: []ServerConfig{{Address: "xxxx"}},
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			conf: NacosConf{
				ServerConfigs: []ServerConfig{{Address: "127.0.0.1:xxx"}},
			},
			wantErr: true,
		},
		{
			name: "normal",
			conf: NacosConf{
				ServerConfigs: []ServerConfig{{Address: "127.0.0.1:1234"}},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.conf.clientParam()
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNacosConf_clientConfig(t *testing.T) {
	tests := []struct {
		name string
		conf NacosConf
		want constant.ClientConfig
	}{
		{
			name: "public namespace",
			conf: NacosConf{
				NamespaceId: "public",
			},
			want: constant.ClientConfig{
				NamespaceId: "",
			},
		},
		{
			name: "other namespace",
			conf: NacosConf{
				NamespaceId: "test",
			},
			want: constant.ClientConfig{
				NamespaceId: "test",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want.NamespaceId, test.conf.clientConfig().NamespaceId)
		})
	}
}
