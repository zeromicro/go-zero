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
		node.CloseX()
		// close again to make sure it's safe
		assert.NotPanics(t, func() {
			node.CloseX()
		})
	})

	t.Run("test blocking node with cluster", func(t *testing.T) {
		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()

		node, err := CreateBlockingNode(New(r.Addr(), Cluster(), WithTLS()))
		assert.NoError(t, err)
		node.CloseX()
		assert.NotPanics(t, func() {
			node.CloseX()
		})
	})

	t.Run("test blocking node with bad type", func(t *testing.T) {
		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()

		_, err = CreateBlockingNode(New(r.Addr(), badType()))
		assert.Error(t, err)
	})
}
