package kube

import (
	"sync"

	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/logx"
	v1 "k8s.io/api/core/v1"
)

type EventHandler struct {
	update    func([]string)
	endpoints map[string]lang.PlaceholderType
	lock      sync.Mutex
}

func NewEventHandler(update func([]string)) *EventHandler {
	return &EventHandler{
		update:    update,
		endpoints: make(map[string]lang.PlaceholderType),
	}
}

func (h *EventHandler) OnAdd(obj interface{}) {
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

func (h *EventHandler) OnDelete(obj interface{}) {
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

func (h *EventHandler) OnUpdate(oldObj, newObj interface{}) {
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
	var targets []string

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
