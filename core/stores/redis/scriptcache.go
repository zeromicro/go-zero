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
	Map map[string]string

	ScriptCache struct {
		atomic.Value
	}
)

func GetScriptCache() *ScriptCache {
	once.Do(func() {
		instance = &ScriptCache{}
		instance.Store(make(Map))
	})

	return instance
}

func (sc *ScriptCache) GetSha(script string) (string, bool) {
	cache := sc.Load().(Map)
	ret, ok := cache[script]
	return ret, ok
}

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
