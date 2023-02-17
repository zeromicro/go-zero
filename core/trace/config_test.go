package trace

import (
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestConfig_parseEndpoint(t *testing.T) {
	logx.Disable()

	c1 := Config{
		Name:     "not UDP",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  kindJaegerUdp,
	}
	c2 := Config{
		Name:     "UDP",
		Endpoint: "localhost:6831",
		Batcher:  kindJaegerUdp,
	}
	host1, port1 := c1.parseEndpoint()
	assert.NotEqual(t, "localhost", host1)
	assert.NotEqual(t, "14268", port1)
	host2, port2 := c2.parseEndpoint()
	assert.Equal(t, "localhost", host2)
	assert.Equal(t, "6831", port2)
}
