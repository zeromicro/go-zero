package consul

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConf(t *testing.T) {
	conf := &Conf{
		Host:     "127.0.0.1:8500",
		ListenOn: "192.168.5.216:9100",
		Key:      "core.rpc",
		Token:    "",
		Tag:      []string{"core", "rpc"},
		Meta:     map[string]string{"Protocol": "grpc"},
		TTL:      0,
	}

	_, err := conf.NewClient()
	assert.Nil(t, err)
}
