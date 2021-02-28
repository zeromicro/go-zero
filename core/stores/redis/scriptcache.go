package redis

import (
	"sync"
	"sync/atomic"
)

var (
	once     sync.Once
	lock     sync.Mutex
	instance *ScriptCache
)

type (
	// Map is an alias of map[string]string.
	Map map[string]string

	// A ScriptCache is a cache that stores a script with its sha key.
	ScriptCache struct {
		atomic.Value
	}
)

// GetScriptCache returns a ScriptCache.
func GetScriptCache() *ScriptCache {
	once.Do(func() {
		instance = &ScriptCache{}
		instance.Store(make(Map))
	})

	return instance
}

// GetSha returns the sha string of given script.
func (sc *ScriptCache) GetSha(script string) (string, bool) {
	cache := sc.Load().(Map)
	ret, ok := cache[script]
	return ret, ok
}

// SetSha sets script with sha into the ScriptCache.
func (sc *ScriptCache) SetSha(script, sha string) {
	lock.Lock()
	defer lock.Unlock()

	cache := sc.Load().(Map)
	newCache := make(Map)
	for k, v := range cache {
		newCache[k] = v
	}
	newCache[script] = sha
	sc.Store(newCache)
}
