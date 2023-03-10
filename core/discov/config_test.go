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

func TestEtcdConf_HasServerID(t *testing.T) {
	tests := []struct {
		EtcdConf
		hasServerID bool
	}{
		{
			EtcdConf: EtcdConf{
				Hosts:    []string{"any"},
				ServerID: -1,
			},
			hasServerID: false,
		},
		{
			EtcdConf: EtcdConf{
				Hosts:    []string{"any"},
				ServerID: 0,
			},
			hasServerID: false,
		},
		{
			EtcdConf: EtcdConf{
				Hosts:    []string{"any"},
				ServerID: 10000,
			},
			hasServerID: true,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.hasServerID, test.EtcdConf.HasServerID())
	}
}
