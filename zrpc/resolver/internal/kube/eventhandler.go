package kube

import (
	"sync"

	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/api/discovery/v1"
	"k8s.io/client-go/tools/cache"
)

var _ cache.ResourceEventHandler = (*EventHandler)(nil)

// EventHandler is ResourceEventHandler implementation.
type EventHandler struct {
	update         func([]string)
	endpointSlices map[string]map[string]lang.PlaceholderType
	lock           sync.Mutex
}

// NewEventHandler returns an EventHandler.
func NewEventHandler(update func([]string)) *EventHandler {
	return &EventHandler{
		update:         update,
		endpointSlices: make(map[string]map[string]lang.PlaceholderType),
	}
}

// OnAdd handles the endpoints add events.
func (h *EventHandler) OnAdd(obj any, _ bool) {
	endpoints, ok := parseEndpointSlice(obj)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.EndpointSlice", obj)
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	old := h.endpoints()
	h.updateEndpointSlice(endpoints)
	if diff(old, h.endpoints()) {
		h.notify()
	}
}

// OnDelete handles the endpoints delete events.
func (h *EventHandler) OnDelete(obj any) {
	endpoints, ok := parseEndpointSlice(obj)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.EndpointSlice", obj)
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	old := h.endpoints()
	if key := endpointSliceKey(endpoints); len(key) > 0 {
		delete(h.endpointSlices, key)
	} else {
		h.deleteEndpointAddresses(endpoints)
	}

	if diff(old, h.endpoints()) {
		h.notify()
	}
}

// OnUpdate handles the endpoints update events.
func (h *EventHandler) OnUpdate(oldObj, newObj any) {
	oldEndpointSlices, ok := parseEndpointSlice(oldObj)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.EndpointSlice", oldObj)
		return
	}

	newEndpointSlices, ok := parseEndpointSlice(newObj)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.EndpointSlice", newObj)
		return
	}

	if oldEndpointSlices.ResourceVersion == newEndpointSlices.ResourceVersion {
		return
	}

	h.Update(newEndpointSlices)
}

// Update updates the endpoints.
func (h *EventHandler) Update(endpoints *v1.EndpointSlice) {
	h.lock.Lock()
	defer h.lock.Unlock()

	old := h.endpoints()
	h.updateEndpointSlice(endpoints)
	if diff(old, h.endpoints()) {
		h.notify()
	}
}

func (h *EventHandler) updateEndpointSlice(endpoints *v1.EndpointSlice) {
	h.endpointSlices[endpointSliceKey(endpoints)] = endpointAddresses(endpoints)
}

func (h *EventHandler) deleteEndpointAddresses(endpoints *v1.EndpointSlice) {
	for _, point := range endpoints.Endpoints {
		for _, address := range point.Addresses {
			for key, addresses := range h.endpointSlices {
				delete(addresses, address)
				if len(addresses) == 0 {
					delete(h.endpointSlices, key)
				}
			}
		}
	}
}

func (h *EventHandler) endpoints() map[string]lang.PlaceholderType {
	endpoints := make(map[string]lang.PlaceholderType)
	for _, slice := range h.endpointSlices {
		for address := range slice {
			endpoints[address] = lang.Placeholder
		}
	}

	return endpoints
}

func (h *EventHandler) notify() {
	endpoints := h.endpoints()
	targets := make([]string, 0, len(endpoints))

	for k := range endpoints {
		targets = append(targets, k)
	}

	h.update(targets)
}

func parseEndpointSlice(obj any) (*v1.EndpointSlice, bool) {
	switch endpoints := obj.(type) {
	case *v1.EndpointSlice:
		return endpoints, true
	case cache.DeletedFinalStateUnknown:
		return parseEndpointSlice(endpoints.Obj)
	default:
		return nil, false
	}
}

func endpointAddresses(endpoints *v1.EndpointSlice) map[string]lang.PlaceholderType {
	addresses := make(map[string]lang.PlaceholderType)
	for _, point := range endpoints.Endpoints {
		for _, address := range point.Addresses {
			addresses[address] = lang.Placeholder
		}
	}

	return addresses
}

func endpointSliceKey(endpoints *v1.EndpointSlice) string {
	if len(endpoints.Namespace) > 0 || len(endpoints.Name) > 0 {
		return endpoints.Namespace + "/" + endpoints.Name
	}

	return string(endpoints.UID)
}

func diff(o, n map[string]lang.PlaceholderType) bool {
	if len(o) != len(n) {
		return true
	}

	for k := range o {
		if _, ok := n[k]; !ok {
			return true
		}
	}

	return false
}
