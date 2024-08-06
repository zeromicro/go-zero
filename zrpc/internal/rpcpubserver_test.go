package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/netx"
)

func TestNewRpcPubServer(t *testing.T) {
	s, err := NewRpcPubServer(discov.EtcdConf{
		User: "user",
		Pass: "pass",
		ID:   10,
	}, "")
	assert.NoError(t, err)
	assert.NotPanics(t, func() {
		s.Start(nil)
	})
}

func TestFigureOutListenOn(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  "192.168.0.5:1234",
			expect: "192.168.0.5:1234",
		},
		{
			input:  "0.0.0.0:8080",
			expect: netx.InternalIp() + ":8080",
		},
		{
			input:  ":8080",
			expect: netx.InternalIp() + ":8080",
		},
		{
			input:  "",
			expect: netx.InternalIp(),
		},
	}

	for _, test := range tests {
		val := figureOutListenOn(test.input)
		assert.Equal(t, test.expect, val)
	}
}
