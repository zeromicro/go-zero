package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestBlockingNode(t *testing.T) {
	r, err := miniredis.Run()
	assert.Nil(t, err)
	node, err := CreateBlockingNode(NewRedis(r.Addr(), NodeType, 0))
	assert.Nil(t, err)
	node.Close()
	node, err = CreateBlockingNode(NewRedis(r.Addr(), ClusterType, 0))
	assert.Nil(t, err)
	node.Close()
}
