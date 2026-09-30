package collection

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/syncx"
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
	// CacheOption defines the method to customize a Cache.
	CacheOption func(cache *Cache)

	// A Cache object is an in-memory cache.
	Cache struct {
		name           string
		lock           sync.Mutex
		data           map[string]any
		expire         time.Duration
		timingWheel    *TimingWheel
		lruCache       lru
		barrier        syncx.SingleFlight
		unstableExpiry mathx.Unstable
		stats          *cacheStat
		// generation counts the registrations of each key in timingWheel. The
		// expiry callback is invoked asynchronously, so a timer may fire and
		// reach Del after the key has been written again. Recording the
		// generation lets Del drop only the timer it was scheduled for.
		generation map[string]uint64
	}
)

// NewCache returns a Cache with given expire.
func NewCache(expire time.Duration, opts ...CacheOption) (*Cache, error) {
	cache := &Cache{
		data:           make(map[string]any),
		expire:         expire,
		lruCache:       emptyLruCache,
		barrier:        syncx.NewSingleFlight(),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		generation:     make(map[string]uint64),
	}

	for _, opt := range opts {
		opt(cache)
	}

	if len(cache.name) == 0 {
		cache.name = defaultCacheName
	}
	cache.stats = newCacheStat(cache.name, cache.size)

	timingWheel, err := NewTimingWheel(time.Second, slots, func(k, v any) {
		key, ok := k.(string)
		if !ok {
			return
		}

		cache.expireKey(key, v)
	})
	if err != nil {
		return nil, err
	}

	cache.timingWheel = timingWheel
	return cache, nil
}

// expireKey removes key on behalf of an expiry that has fired. The value the
// timing wheel was given carries the generation of the write that scheduled
// this timer, so an expiry belonging to a superseded write is ignored rather
// than deleting the key and cancelling the timer of the write that replaced it.
func (c *Cache) expireKey(key string, v any) {
	gen, ok := v.(uint64)
	if !ok {
		c.Del(key)
		return
	}

	c.lock.Lock()
	if c.generation[key] != gen {
		// The key was written again after this timer was scheduled.
		// The newer timer will expire it in turn.
		c.lock.Unlock()
		return
	}
	delete(c.data, key)
	c.lruCache.remove(key)
	c.lock.Unlock()

	c.timingWheel.RemoveTimer(key)
}

// Del deletes the item with the given key from c.
func (c *Cache) Del(key string) {
	c.lock.Lock()
	delete(c.data, key)
	delete(c.generation, key)
	c.lruCache.remove(key)
	c.lock.Unlock()

	// RemoveTimer is called outside the lock to avoid performance impact from this
	// potentially time-consuming operation. Data integrity is maintained by lruCache,
	// which will eventually evict any remaining entries when capacity is exceeded.
	c.timingWheel.RemoveTimer(key)
}

// Get returns the item with the given key from c.
func (c *Cache) Get(key string) (any, bool) {
	value, ok := c.doGet(key)
	if ok {
		c.stats.IncrementHit()
	} else {
		c.stats.IncrementMiss()
	}

	return value, ok
}

// Set sets value into c with key.
func (c *Cache) Set(key string, value any) {
	c.SetWithExpire(key, value, c.expire)
}

// SetWithExpire sets value into c with key and expire with the given value.
func (c *Cache) SetWithExpire(key string, value any, expire time.Duration) {
	c.lock.Lock()
	c.data[key] = value
	c.lruCache.add(key)
	// Record the registration so that an expiry firing for an earlier write
	// of this key can tell it has been superseded. The timing wheel is the only
	// channel back to the callback, so the generation is what it carries.
	c.generation[key]++
	gen := c.generation[key]
	c.lock.Unlock()

	// MoveTimer carries no value, so the key is re-armed with SetTimer in order
	// to hand the new generation to the wheel. SetTimer on an existing key
	// updates the value and reschedules it, which is what an update needs.
	c.timingWheel.SetTimer(key, gen, c.unstableExpiry.AroundDuration(expire))
}

// Take returns the item with the given key.
// If the item is in c, return it directly.
// If not, use fetch method to get the item, set into c and return it.
func (c *Cache) Take(key string, fetch func() (any, error)) (any, error) {
	if val, ok := c.doGet(key); ok {
		c.stats.IncrementHit()
		return val, nil
	}

	var fresh bool
	val, err := c.barrier.Do(key, func() (any, error) {
		// because O(1) on map search in memory, and fetch is an IO query,
		// so we do double-check, cache might be taken by another call
		if val, ok := c.doGet(key); ok {
			return val, nil
		}

		v, e := fetch()
		if e != nil {
			return nil, e
		}

		fresh = true
		c.Set(key, v)
		return v, nil
	})
	if err != nil {
		return nil, err
	}

	if fresh {
		c.stats.IncrementMiss()
		return val, nil
	}

	// got the result from previous ongoing query
	c.stats.IncrementHit()
	return val, nil
}

func (c *Cache) doGet(key string) (any, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.data[key]
	if ok {
		c.lruCache.add(key)
	}

	return value, ok
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

// WithLimit customizes a Cache with items up to limit.
func WithLimit(limit int) CacheOption {
	return func(cache *Cache) {
		if limit > 0 {
			cache.lruCache = newKeyLru(limit, cache.onEvict)
		}
	}
}

// WithName customizes a Cache with the given name.
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

	for range ticker.C {
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
