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
