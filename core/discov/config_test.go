package discov

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		EtcdConf
		pass bool
	}{
		{
			EtcdConf: EtcdConf{},
			pass:     false,
		},
		{
			EtcdConf: EtcdConf{
				Key: "any",
			},
			pass: false,
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
			},
			pass: false,
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				Key:   "key",
			},
			pass: true,
		},
	}

	for _, test := range tests {
		if test.pass {
			assert.Nil(t, test.EtcdConf.Validate())
		} else {
			assert.NotNil(t, test.EtcdConf.Validate())
		}
	}
}

func TestEtcdConf_HasAccount(t *testing.T) {
	tests := []struct {
		EtcdConf
		hasAccount bool
	}{
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				Key:   "key",
			},
			hasAccount: false,
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				Key:   "key",
				User:  "foo",
			},
			hasAccount: false,
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				Key:   "key",
				User:  "foo",
				Pass:  "bar",
			},
			hasAccount: true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.hasAccount, test.EtcdConf.HasAccount())
	}
}

func TestEtcdConf_HasID(t *testing.T) {
	tests := []struct {
		EtcdConf
		hasServerID bool
	}{
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				ID:    -1,
			},
			hasServerID: false,
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				ID:    0,
			},
			hasServerID: false,
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				ID:    10000,
			},
			hasServerID: true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.hasServerID, test.EtcdConf.HasID())
	}
}

func TestEtcdConf_HasTLS(t *testing.T) {
	tests := []struct {
		name string
		conf EtcdConf
		want bool
	}{
		{
			name: "empty config",
			conf: EtcdConf{},
			want: false,
		},
		{
			name: "missing CertFile",
			conf: EtcdConf{
				CertKeyFile: "key",
				CACertFile:  "ca",
			},
			want: false,
		},
		{
			name: "missing CertKeyFile",
			conf: EtcdConf{
				CertFile:   "cert",
				CACertFile: "ca",
			},
			want: false,
		},
		{
			name: "missing CACertFile",
			conf: EtcdConf{
				CertFile:    "cert",
				CertKeyFile: "key",
			},
			want: false,
		},
		{
			name: "valid config",
			conf: EtcdConf{
				CertFile:    "cert",
				CertKeyFile: "key",
				CACertFile:  "ca",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.conf.HasTLS()
			assert.Equal(t, tt.want, got)
		})
	}
}
