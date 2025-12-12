package discov

import "github.com/zeromicro/go-zero/core/discov/internal"

// Container is the interface for handling add and delete events of key-value pairs.
type Container interface {
	// OnAdd is called when a new key-value pair is added.
	OnAdd(kv internal.KV)
	// OnDelete is called when a key-value pair is deleted.
	OnDelete(kv internal.KV)
	// addListener adds a listener function that will be called on changes.
	addListener(listener func())
	// getValues returns all the current values.
	getValues() []string
	// notifyChange notifies all listeners of a change.
	notifyChange()
}
