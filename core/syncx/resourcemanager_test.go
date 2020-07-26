package syncx

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyResource struct {
	age int
}

func (dr *dummyResource) Close() error {
	return errors.New("close")
}

func TestResourceManager_GetResource(t *testing.T) {
	manager := NewResourceManager()
	defer manager.Close()

	var age int
	for i := 0; i < 10; i++ {
		val, err := manager.GetResource("key", func() (io.Closer, error) {
			age++
			return &dummyResource{
				age: age,
			}, nil
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, val.(*dummyResource).age)
	}
}

func TestResourceManager_GetResourceError(t *testing.T) {
	manager := NewResourceManager()
	defer manager.Close()

	for i := 0; i < 10; i++ {
		_, err := manager.GetResource("key", func() (io.Closer, error) {
			return nil, errors.New("fail")
		})
		assert.NotNil(t, err)
	}
}
