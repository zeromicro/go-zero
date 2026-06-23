package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDirectTarget(t *testing.T) {
	target := BuildDirectTarget([]string{"localhost:123", "localhost:456"})
	assert.Equal(t, "direct:///localhost:123,localhost:456", target)
}

func TestBuildDiscovTarget(t *testing.T) {
	target := BuildDiscovTarget([]string{"localhost:123", "localhost:456"}, "foo")
	assert.Equal(t, "etcd:///localhost:123,localhost:456?key=foo", target)
}

func TestBuildDiscovTargetWithSlashKey(t *testing.T) {
	target := BuildDiscovTarget([]string{"localhost:2379"}, "/grpc/my-service")
	assert.Equal(t, "etcd:///localhost:2379?key=%2Fgrpc%2Fmy-service", target)
}
