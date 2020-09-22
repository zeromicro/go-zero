package internal

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
	assert.Equal(t, "discov://localhost:123,localhost:456/foo", target)
}
