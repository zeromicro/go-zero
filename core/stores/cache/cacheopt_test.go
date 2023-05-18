package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheOptions(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		o := newOptions()
		assert.Equal(t, defaultExpiry, o.Expiry)
		assert.Equal(t, defaultNotFoundExpiry, o.NotFoundExpiry)
	})

	t.Run("with expiry", func(t *testing.T) {
		o := newOptions(WithExpiry(time.Second))
		assert.Equal(t, time.Second, o.Expiry)
		assert.Equal(t, defaultNotFoundExpiry, o.NotFoundExpiry)
	})

	t.Run("with not found expiry", func(t *testing.T) {
		o := newOptions(WithNotFoundExpiry(time.Second))
		assert.Equal(t, defaultExpiry, o.Expiry)
		assert.Equal(t, time.Second, o.NotFoundExpiry)
	})
}
