package redis

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestBlockingNode(t *testing.T) {
	t.Run("test blocking node", func(t *testing.T) {
		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()

		node, err := CreateBlockingNode(New(r.Addr()))
		assert.NoError(t, err)
		node.Close()
		// close again to make sure it's safe
		assert.NotPanics(t, func() {
			node.Close()
		})
	})

	t.Run("test blocking node with cluster", func(t *testing.T) {
		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()

		node, err := CreateBlockingNode(New(r.Addr(), Cluster(), WithTLS()))
		assert.NoError(t, err)
		node.Close()
		assert.NotPanics(t, func() {
			node.Close()
		})
	})

	t.Run("test blocking node with bad type", func(t *testing.T) {
		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()

		_, err = CreateBlockingNode(New(r.Addr(), badType()))
		assert.Error(t, err)
	})

	t.Run("test blocking node with protocol and identity", func(t *testing.T) {
		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()

		node, err := CreateBlockingNode(New(r.Addr(), WithProtocol(2), WithIdentity()))
		assert.NoError(t, err)
		bridge, ok := node.(*clientBridge)
		assert.True(t, ok)
		assert.Equal(t, 2, bridge.Options().Protocol)
		assert.True(t, bridge.Options().DisableIdentity)
		node.Close()
	})

	t.Run("test blocking node with cluster, protocol and identity", func(t *testing.T) {
		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()

		node, err := CreateBlockingNode(New(r.Addr(), Cluster(), WithProtocol(2), WithIdentity()))
		assert.NoError(t, err)
		bridge, ok := node.(*clusterBridge)
		assert.True(t, ok)
		assert.Equal(t, 2, bridge.Options().Protocol)
		assert.True(t, bridge.Options().DisableIdentity)
		node.Close()
	})
}
