package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefCleaner(t *testing.T) {
	var count int
	clean := func() {
		count += 1
	}

	cleaner := NewRefResource(clean)
	err := cleaner.Use()
	assert.Nil(t, err)
	err = cleaner.Use()
	assert.Nil(t, err)
	cleaner.Clean()
	cleaner.Clean()
	assert.Equal(t, 1, count)
	cleaner.Clean()
	cleaner.Clean()
	assert.Equal(t, 1, count)
	assert.Equal(t, ErrUseOfCleaned, cleaner.Use())
}
