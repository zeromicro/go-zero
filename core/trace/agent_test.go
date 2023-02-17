package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestStartAgent(t *testing.T) {
	logx.Disable()

	const (
		endpoint1  = "localhost:1234"
		endpoint2  = "remotehost:1234"
		endpoint3  = "localhost:1235"
		endpoint4  = "localhost:1236"
		agentHost1 = "localhost"
		agentPort1 = "6831"
	)
	c1 := Config{
		Name: "foo",
	}
	c2 := Config{
		Name:     "bar",
		Endpoint: endpoint1,
		Batcher:  kindJaeger,
	}
	c3 := Config{
		Name:     "any",
		Endpoint: endpoint2,
		Batcher:  kindZipkin,
	}
	c4 := Config{
		Name:     "bla",
		Endpoint: endpoint3,
		Batcher:  "otlp",
	}
	c5 := Config{
		Name:     "grpc",
		Endpoint: endpoint3,
		Batcher:  kindOtlpGrpc,
	}
	c6 := Config{
		Name:     "otlphttp",
		Endpoint: endpoint4,
		Batcher:  kindOtlpHttp,
	}
	c7 := Config{
		Name:      "jaegerUDP",
		AgentHost: agentHost1,
		AgentPort: agentPort1,
		Batcher:   kindJaeger,
	}
	c8 := Config{
		Name:      "jaegerUDP",
		AgentHost: agentHost1,
		AgentPort: agentPort1,
		Endpoint:  endpoint1,
		Batcher:   kindJaeger,
	}

	StartAgent(c1)
	StartAgent(c1)
	StartAgent(c2)
	StartAgent(c3)
	StartAgent(c4)
	StartAgent(c5)
	StartAgent(c6)
	StartAgent(c7)
	StartAgent(c8)

	lock.Lock()
	defer lock.Unlock()

	// because remotehost cannot be resolved
	assert.Equal(t, 5, len(agents))
	_, ok := agents[""]
	assert.True(t, ok)
	_, ok = agents[endpoint1]
	assert.True(t, ok)
	_, ok = agents[endpoint2]
	assert.False(t, ok)
	_, ok = agents[c2.getEndpoint()]
	assert.True(t, ok)
	_, ok = agents[c3.getEndpoint()]
	assert.False(t, ok)
	_, ok = agents[c7.getEndpoint()]
	assert.True(t, ok)
	_, ok = agents[c8.getEndpoint()]
	assert.True(t, ok)
}
