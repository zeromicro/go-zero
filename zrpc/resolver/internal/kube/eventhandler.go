package kube

import (
	"sync"

	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

var _ cache.ResourceEventHandler = (*EventHandler)(nil)

// EventHandler is ResourceEventHandler implementation.
type EventHandler struct {
	update    func([]string)
	endpoints map[string]lang.PlaceholderType
	lock      sync.Mutex
}

// NewEventHandler returns an EventHandler.
func NewEventHandler(update func([]string)) *EventHandler {
	return &EventHandler{
		update:    update,
		endpoints: make(map[string]lang.PlaceholderType),
	}
}

// OnAdd handles the endpoints add events.
func (h *EventHandler) OnAdd(obj any, _ bool) {
	endpoints, ok := obj.(*v1.Endpoints)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.Endpoints", obj)
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	var changed bool
	for _, sub := range endpoints.Subsets {
		for _, point := range sub.Addresses {
			if _, ok := h.endpoints[point.IP]; !ok {
				h.endpoints[point.IP] = lang.Placeholder
				changed = true
			}
		}
	}

	if changed {
		h.notify()
	}
}

// OnDelete handles the endpoints delete events.
func (h *EventHandler) OnDelete(obj any) {
	endpoints, ok := obj.(*v1.Endpoints)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.Endpoints", obj)
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	var changed bool
	for _, sub := range endpoints.Subsets {
		for _, point := range sub.Addresses {
			if _, ok := h.endpoints[point.IP]; ok {
				delete(h.endpoints, point.IP)
				changed = true
			}
		}
	}

	if changed {
		h.notify()
	}
}

// OnUpdate handles the endpoints update events.
func (h *EventHandler) OnUpdate(oldObj, newObj any) {
	oldEndpoints, ok := oldObj.(*v1.Endpoints)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.Endpoints", oldObj)
		return
	}

	newEndpoints, ok := newObj.(*v1.Endpoints)
	if !ok {
		logx.Errorf("%v is not an object with type *v1.Endpoints", newObj)
		return
	}

	if oldEndpoints.ResourceVersion == newEndpoints.ResourceVersion {
		return
	}

	h.Update(newEndpoints)
}

// Update updates the endpoints.
func (h *EventHandler) Update(endpoints *v1.Endpoints) {
	h.lock.Lock()
	defer h.lock.Unlock()

	old := h.endpoints
	h.endpoints = make(map[string]lang.PlaceholderType)
	for _, sub := range endpoints.Subsets {
		for _, point := range sub.Addresses {
			h.endpoints[point.IP] = lang.Placeholder
		}
	}

	if diff(old, h.endpoints) {
		h.notify()
	}
}

func (h *EventHandler) notify() {
	targets := make([]string, 0, len(h.endpoints))

	for k := range h.endpoints {
		targets = append(targets, k)
	}

	h.update(targets)
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
