package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartAgent(t *testing.T) {
	const (
		endpoint1 = "localhost:1234"
		endpoint2 = "remotehost:1234"
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

	StartAgent(c1)
	StartAgent(c1)
	StartAgent(c2)
	StartAgent(c3)

	lock.Lock()
	defer lock.Unlock()

	assert.Equal(t, 3, len(agents))
	_, ok := agents[""]
	assert.True(t, ok)
	_, ok = agents[endpoint1]
	assert.True(t, ok)
	_, ok = agents[endpoint2]
	assert.True(t, ok)
}
