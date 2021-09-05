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
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				Key:   "key",
				User:  "root",
				Pass:  "pass",
			},
			pass: true,
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				Key:   "key",
				User:  "root",
				Pass:  "",
			},
			pass: false,
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				Key:   "key",
				User:  "",
				Pass:  "pass",
			},
			pass: false,
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				Key:   "key",
				User:  "",
				Pass:  "",
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

func TestEnableAuth(t *testing.T) {
	tests := []struct {
		EtcdConf
		pass bool
	}{
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				User:  "root",
				Pass:  "password",
			},
			pass: true,
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				User:  "",
				Pass:  "password",
			},
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				User:  "root",
				Pass:  "",
			},
		},
		{
			EtcdConf: EtcdConf{
				Hosts: []string{"any"},
				User:  "",
				Pass:  "",
			},
		},
	}

	for _, test := range tests {
		if test.pass {
			assert.True(t, test.EtcdConf.EnableAuth())
		} else {
			assert.False(t, test.EtcdConf.EnableAuth())
		}
	}
}
