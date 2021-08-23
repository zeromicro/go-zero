package syncx

import (
	"io"
	"sync"

	"github.com/tal-tech/go-zero/core/errorx"
)

// A ResourceManager is a manager that used to manage resources.
type ResourceManager struct {
	resources   map[string]io.Closer
	sharedCalls SharedCalls
	lock        sync.RWMutex
}

// NewResourceManager returns a ResourceManager.
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resources:   make(map[string]io.Closer),
		sharedCalls: NewSharedCalls(),
	}
}

// Close closes the manager.
func (manager *ResourceManager) Close() error {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	var be errorx.BatchError
	for _, resource := range manager.resources {
		if err := resource.Close(); err != nil {
			be.Add(err)
		}
	}

	return be.Err()
}

// GetResource returns the resource associated with given key.
func (manager *ResourceManager) GetResource(key string, create func() (io.Closer, error)) (io.Closer, error) {
	val, err := manager.sharedCalls.Do(key, func() (interface{}, error) {
		manager.lock.RLock()
		resource, ok := manager.resources[key]
		manager.lock.RUnlock()
		if ok {
			return resource, nil
		}

		resource, err := create()
		if err != nil {
			return nil, err
		}

		manager.lock.Lock()
		manager.resources[key] = resource
		manager.lock.Unlock()

		return resource, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(io.Closer), nil
}
