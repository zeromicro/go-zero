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
	update    func([]string)
	endpoints map[string]lang.PlaceholderType
	slices    map[string]map[string]lang.PlaceholderType
	lock      sync.Mutex
}

// NewEventHandler returns an EventHandler.
func NewEventHandler(update func([]string)) *EventHandler {
	return &EventHandler{
		update:    update,
		endpoints: make(map[string]lang.PlaceholderType),
		slices:    make(map[string]map[string]lang.PlaceholderType),
	}
}

// OnAdd handles the endpoints add events.
func (h *EventHandler) OnAdd(obj any, _ bool) {
	endpoints, ok := obj.(*v1.EndpointSlice)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.EndpointSlice", obj)
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	h.slices[sliceKey(endpoints)] = readyAddresses(endpoints)

	h.rebuildEndpoints()
}

// OnDelete handles the endpoints delete events.
func (h *EventHandler) OnDelete(obj any) {
	endpoints, ok := obj.(*v1.EndpointSlice)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.EndpointSlice", obj)
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	delete(h.slices, sliceKey(endpoints))

	h.rebuildEndpoints()
}

// OnUpdate handles the endpoints update events.
func (h *EventHandler) OnUpdate(oldObj, newObj any) {
	oldEndpointSlices, ok := oldObj.(*v1.EndpointSlice)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.EndpointSlice", oldObj)
		return
	}

	newEndpointSlices, ok := newObj.(*v1.EndpointSlice)
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

	h.slices[sliceKey(endpoints)] = readyAddresses(endpoints)
	h.rebuildEndpoints()
}

func (h *EventHandler) rebuildEndpoints() {
	old := h.endpoints
	h.endpoints = make(map[string]lang.PlaceholderType)
	for _, addresses := range h.slices {
		for address := range addresses {
			h.endpoints[address] = lang.Placeholder
		}
	}
	if diff(old, h.endpoints) {
		h.notify()
	}
}

func readyAddresses(endpoints *v1.EndpointSlice) map[string]lang.PlaceholderType {
	addresses := make(map[string]lang.PlaceholderType)
	for _, point := range endpoints.Endpoints {
		if isReady(point) {
			for _, address := range point.Addresses {
				addresses[address] = lang.Placeholder
			}
		}
	}
	return addresses
}

func sliceKey(endpoints *v1.EndpointSlice) string {
	return endpoints.Namespace + "/" + endpoints.Name
}

func (h *EventHandler) notify() {
	targets := make([]string, 0, len(h.endpoints))

	for k := range h.endpoints {
		targets = append(targets, k)
	}

	h.update(targets)
}

// isReady reports whether the endpoint should receive new traffic, the same
// set of addresses the Endpoints API listed as ready. Kubernetes reports
// terminating endpoints as not ready, and a nil Ready condition means unknown,
// which consumers should treat as ready.
func isReady(point v1.Endpoint) bool {
	return point.Conditions.Ready == nil || *point.Conditions.Ready
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
