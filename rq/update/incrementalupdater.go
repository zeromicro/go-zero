package update

import (
	"sync"
	"time"

	"zero/core/hash"
	"zero/core/stringx"
)

const (
	incrementalStep = 5
	stepDuration    = time.Second * 3
)

type (
	updateEvent struct {
		keys    []string
		newKey  string
		servers []string
	}

	UpdateFunc func(change ServerChange)

	IncrementalUpdater struct {
		lock          sync.Mutex
		started       bool
		taskChan      chan updateEvent
		updates       ServerChange
		updateFn      UpdateFunc
		pendingEvents []updateEvent
	}
)

func NewIncrementalUpdater(updateFn UpdateFunc) *IncrementalUpdater {
	return &IncrementalUpdater{
		taskChan: make(chan updateEvent),
		updates: ServerChange{
			Current: Snapshot{
				Keys:         make([]string, 0),
				WeightedKeys: make([]weightedKey, 0),
			},
			Servers: make([]string, 0),
		},
		updateFn: updateFn,
	}
}

func (ru *IncrementalUpdater) Update(keys []string, servers []string, newKey string) {
	ru.lock.Lock()
	defer ru.lock.Unlock()

	if !ru.started {
		go ru.run()
		ru.started = true
	}

	ru.taskChan <- updateEvent{
		keys:    keys,
		newKey:  newKey,
		servers: servers,
	}
}

// Return true if incremental update is done
func (ru *IncrementalUpdater) advance() bool {
	previous := ru.updates.Current
	keys := make([]string, 0)
	weightedKeys := make([]weightedKey, 0)
	servers := ru.updates.Servers
	for _, key := range ru.updates.Current.Keys {
		keys = append(keys, key)
	}
	for _, wkey := range ru.updates.Current.WeightedKeys {
		weight := wkey.Weight + incrementalStep
		if weight >= hash.TopWeight {
			keys = append(keys, wkey.Key)
		} else {
			weightedKeys = append(weightedKeys, weightedKey{
				Key:    wkey.Key,
				Weight: weight,
			})
		}
	}

	for _, event := range ru.pendingEvents {
		// ignore reload events
		if len(event.newKey) == 0 || len(event.servers) == 0 {
			continue
		}

		// anyway, add the servers, just to avoid missing notify any server
		servers = stringx.Union(servers, event.servers)
		if keyExists(keys, weightedKeys, event.newKey) {
			continue
		}

		weightedKeys = append(weightedKeys, weightedKey{
			Key:    event.newKey,
			Weight: incrementalStep,
		})
	}

	// clear pending events
	ru.pendingEvents = ru.pendingEvents[:0]

	change := ServerChange{
		Previous: previous,
		Current: Snapshot{
			Keys:         keys,
			WeightedKeys: weightedKeys,
		},
		Servers: servers,
	}
	ru.updates = change
	ru.updateFn(change)

	return len(weightedKeys) == 0
}

func (ru *IncrementalUpdater) run() {
	defer func() {
		ru.lock.Lock()
		ru.started = false
		ru.lock.Unlock()
	}()

	ticker := time.NewTicker(stepDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ru.advance() {
				return
			}
		case event := <-ru.taskChan:
			ru.updateKeys(event)
		}
	}
}

func (ru *IncrementalUpdater) updateKeys(event updateEvent) {
	isWeightedKey := func(key string) bool {
		for _, wkey := range ru.updates.Current.WeightedKeys {
			if wkey.Key == key {
				return true
			}
		}

		return false
	}

	keys := make([]string, 0, len(event.keys))
	for _, key := range event.keys {
		if !isWeightedKey(key) {
			keys = append(keys, key)
		}
	}

	ru.updates.Current.Keys = keys
	ru.pendingEvents = append(ru.pendingEvents, event)
}

func keyExists(keys []string, weightedKeys []weightedKey, key string) bool {
	for _, each := range keys {
		if key == each {
			return true
		}
	}

	for _, wkey := range weightedKeys {
		if wkey.Key == key {
			return true
		}
	}

	return false
}
