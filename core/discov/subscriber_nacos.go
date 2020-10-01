package discov

import (
	"sync"
)

type (
	SubscriberNacos struct {
		listeners []func()
		lock      sync.Mutex
	}
)

func NewSubscriberNacos(endpoints []string, key string, opts ...SubOption) (Subscriber, error) {
	var subOpts subOptions
	for _, opt := range opts {
		opt(&subOpts)
	}

	sub := &SubscriberNacos{}

	return sub, nil
}

func (s *SubscriberNacos) AddListener(listener func()) {
	s.lock.Lock()
	s.listeners = append(s.listeners, listener)
	s.lock.Unlock()
}

func (s *SubscriberNacos) Values() []string {
	return []string{}
}
