package syncx

import "sync"

// A ManagedResource is used to manage a resource that might be broken and refetched, like a connection.
type ManagedResource struct {
	resource interface{}
	lock     sync.RWMutex
	generate func() interface{}
	equals   func(a, b interface{}) bool
}

// NewManagedResource returns a ManagedResource.
func NewManagedResource(generate func() interface{}, equals func(a, b interface{}) bool) *ManagedResource {
	return &ManagedResource{
		generate: generate,
		equals:   equals,
	}
}

// MarkBroken marks the resource broken.
func (mr *ManagedResource) MarkBroken(resource interface{}) {
	mr.lock.Lock()
	defer mr.lock.Unlock()

	if mr.equals(mr.resource, resource) {
		mr.resource = nil
	}
}

// Take takes the resource, if not loaded, generates it.
func (mr *ManagedResource) Take() interface{} {
	mr.lock.RLock()
	resource := mr.resource
	mr.lock.RUnlock()

	if resource != nil {
		return resource
	}

	mr.lock.Lock()
	defer mr.lock.Unlock()
	// maybe another Take() call already generated the resource.
	if mr.resource == nil {
		mr.resource = mr.generate()
	}
	return mr.resource
}
