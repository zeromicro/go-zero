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

type trackedResource struct {
	closed bool
}

func (tr *trackedResource) Close() error {
	tr.closed = true
	return nil
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

func TestResourceManager_Close(t *testing.T) {
	manager := NewResourceManager()
	defer manager.Close()

	for i := 0; i < 10; i++ {
		_, err := manager.GetResource("key", func() (io.Closer, error) {
			return nil, errors.New("fail")
		})
		assert.NotNil(t, err)
	}

	if assert.NoError(t, manager.Close()) {
		assert.Equal(t, 0, len(manager.resources))
	}
}

func TestResourceManager_UseAfterClose(t *testing.T) {
	manager := NewResourceManager()
	defer manager.Close()

	_, err := manager.GetResource("key", func() (io.Closer, error) {
		return nil, errors.New("fail")
	})
	assert.NotNil(t, err)
	if assert.NoError(t, manager.Close()) {
		_, err = manager.GetResource("key", func() (io.Closer, error) {
			return nil, errors.New("fail")
		})
		assert.NotNil(t, err)

		assert.Panics(t, func() {
			_, err = manager.GetResource("key", func() (io.Closer, error) {
				return &dummyResource{age: 123}, nil
			})
		})
	}
}

func TestResourceManager_Inject(t *testing.T) {
	manager := NewResourceManager()
	defer manager.Close()

	manager.Inject("key", &dummyResource{
		age: 10,
	})

	val, err := manager.GetResource("key", func() (io.Closer, error) {
		return nil, nil
	})
	assert.Nil(t, err)
	assert.Equal(t, 10, val.(*dummyResource).age)
}

func TestResourceManager_RemoveResource(t *testing.T) {
	manager := NewResourceManager()
	defer manager.Close()

	resource := &trackedResource{}
	manager.Inject("key", resource)

	assert.NoError(t, manager.RemoveResource("key"))
	assert.True(t, resource.closed)

	_, ok := manager.resources["key"]
	assert.False(t, ok)
}

func TestResourceManager_RemoveResourceMissing(t *testing.T) {
	manager := NewResourceManager()
	defer manager.Close()

	assert.NoError(t, manager.RemoveResource("missing"))
}
