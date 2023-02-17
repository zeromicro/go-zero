package trace

import (
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestConfig_getEndpointHost(t *testing.T) {
	logx.Disable()

	c1 := Config{
		Endpoint: "http://localhost:14268/api/traces",
	}
	c2 := Config{
		Endpoint: "localhost:6831",
	}
	assert.NotEqual(t, "localhost", c1.getEndpointHost())
	assert.NotEqual(t, "14268", c1.getEndpointPort())
	assert.Equal(t, "localhost", c2.getEndpointHost())
	assert.Equal(t, "6831", c2.getEndpointPort())
}
