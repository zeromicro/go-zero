package collection

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/mathx"
	"github.com/tal-tech/go-zero/core/syncx"
)

const (
	defaultCacheName = "proc"
	slots            = 300
	statInterval     = time.Minute
	// make the expiry unstable to avoid lots of cached items expire at the same time
	// make the unstable expiry to be [0.95, 1.05] * seconds
	expiryDeviation = 0.05
)

var emptyLruCache = emptyLru{}

type (
	CacheOption func(cache *Cache)

	Cache struct {
		name           string
		lock           sync.Mutex
		data           map[string]interface{}
		evicts         *list.List
		expire         time.Duration
		timingWheel    *TimingWheel
		lruCache       lru
		barrier        syncx.SharedCalls
		unstableExpiry mathx.Unstable
		stats          *cacheStat
	}
)

func NewCache(expire time.Duration, opts ...CacheOption) (*Cache, error) {
	cache := &Cache{
		data:           make(map[string]interface{}),
		expire:         expire,
		lruCache:       emptyLruCache,
		barrier:        syncx.NewSharedCalls(),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
	}

	for _, opt := range opts {
		opt(cache)
	}

	if len(cache.name) == 0 {
		cache.name = defaultCacheName
	}
	cache.stats = newCacheStat(cache.name, cache.size)

	timingWheel, err := NewTimingWheel(time.Second, slots, func(k, v interface{}) {
		key, ok := k.(string)
		if !ok {
			return
		}

		cache.Del(key)
	})
	if err != nil {
		return nil, err
	}

	cache.timingWheel = timingWheel
	return cache, nil
}

func (c *Cache) Del(key string) {
	c.lock.Lock()
	delete(c.data, key)
	c.lruCache.remove(key)
	c.lock.Unlock()
	c.timingWheel.RemoveTimer(key)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.lock.Lock()
	value, ok := c.data[key]
	if ok {
		c.lruCache.add(key)
	}
	c.lock.Unlock()
	if ok {
		c.stats.IncrementHit()
	} else {
		c.stats.IncrementMiss()
	}

	return value, ok
}

func (c *Cache) Set(key string, value interface{}) {
	c.lock.Lock()
	_, ok := c.data[key]
	c.data[key] = value
	c.lruCache.add(key)
	c.lock.Unlock()

	expiry := c.unstableExpiry.AroundDuration(c.expire)
	if ok {
		c.timingWheel.MoveTimer(key, expiry)
	} else {
		c.timingWheel.SetTimer(key, value, expiry)
	}
}

func (c *Cache) Take(key string, fetch func() (interface{}, error)) (interface{}, error) {
	val, fresh, err := c.barrier.DoEx(key, func() (interface{}, error) {
		v, e := fetch()
		if e != nil {
			return nil, e
		}

		c.Set(key, v)
		return v, nil
	})
	if err != nil {
		return nil, err
	}

	if fresh {
		c.stats.IncrementMiss()
		return val, nil
	} else {
		// got the result from previous ongoing query
		c.stats.IncrementHit()
	}

	return val, nil
}

func (c *Cache) onEvict(key string) {
	// already locked
	delete(c.data, key)
	c.timingWheel.RemoveTimer(key)
}

func (c *Cache) size() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

func WithLimit(limit int) CacheOption {
	return func(cache *Cache) {
		if limit > 0 {
			cache.lruCache = newKeyLru(limit, cache.onEvict)
		}
	}
}

func WithName(name string) CacheOption {
	return func(cache *Cache) {
		cache.name = name
	}
}

type (
	lru interface {
		add(key string)
		remove(key string)
	}

	emptyLru struct{}

	keyLru struct {
		limit    int
		evicts   *list.List
		elements map[string]*list.Element
		onEvict  func(key string)
	}
)

func (elru emptyLru) add(string) {
}

func (elru emptyLru) remove(string) {
}

func newKeyLru(limit int, onEvict func(key string)) *keyLru {
	return &keyLru{
		limit:    limit,
		evicts:   list.New(),
		elements: make(map[string]*list.Element),
		onEvict:  onEvict,
	}
}

func (klru *keyLru) add(key string) {
	if elem, ok := klru.elements[key]; ok {
		klru.evicts.MoveToFront(elem)
		return
	}

	// Add new item
	elem := klru.evicts.PushFront(key)
	klru.elements[key] = elem

	// Verify size not exceeded
	if klru.evicts.Len() > klru.limit {
		klru.removeOldest()
	}
}

func (klru *keyLru) remove(key string) {
	if elem, ok := klru.elements[key]; ok {
		klru.removeElement(elem)
	}
}

func (klru *keyLru) removeOldest() {
	elem := klru.evicts.Back()
	if elem != nil {
		klru.removeElement(elem)
	}
}

func (klru *keyLru) removeElement(e *list.Element) {
	klru.evicts.Remove(e)
	key := e.Value.(string)
	delete(klru.elements, key)
	klru.onEvict(key)
}

type cacheStat struct {
	name         string
	hit          uint64
	miss         uint64
	sizeCallback func() int
}

func newCacheStat(name string, sizeCallback func() int) *cacheStat {
	st := &cacheStat{
		name:         name,
		sizeCallback: sizeCallback,
	}
	go st.statLoop()
	return st
}

func (cs *cacheStat) IncrementHit() {
	atomic.AddUint64(&cs.hit, 1)
}

func (cs *cacheStat) IncrementMiss() {
	atomic.AddUint64(&cs.miss, 1)
}

func (cs *cacheStat) statLoop() {
	ticker := time.NewTicker(statInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hit := atomic.SwapUint64(&cs.hit, 0)
			miss := atomic.SwapUint64(&cs.miss, 0)
			total := hit + miss
			if total == 0 {
				continue
			}
			percent := 100 * float32(hit) / float32(total)
			logx.Statf("cache(%s) - qpm: %d, hit_ratio: %.1f%%, elements: %d, hit: %d, miss: %d",
				cs.name, total, percent, cs.sizeCallback(), hit, miss)
		}
	}
}
