package syncx

import "sync"

type ManagedResource struct {
	resource interface{}
	lock     sync.RWMutex
	generate func() interface{}
	equals   func(a, b interface{}) bool
}

func NewManagedResource(generate func() interface{}, equals func(a, b interface{}) bool) *ManagedResource {
	return &ManagedResource{
		generate: generate,
		equals:   equals,
	}
}

func (mr *ManagedResource) MarkBroken(resource interface{}) {
	mr.lock.Lock()
	defer mr.lock.Unlock()

	if mr.equals(mr.resource, resource) {
		mr.resource = nil
	}
}

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
