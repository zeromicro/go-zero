# 通过 collection.Cache 进行缓存

go-zero微服务框架中提供了许多开箱即用的工具，好的工具不仅能提升服务的性能而且还能提升代码的鲁棒性避免出错，实现代码风格的统一方便他人阅读等等，本系列文章将分别介绍go-zero框架中工具的使用及其实现原理  

## 进程内缓存工具[collection.Cache](https://github.com/tal-tech/go-zero/tree/master/core/collection/cache.go)

在做服务器开发的时候，相信都会遇到使用缓存的情况，go-zero 提供的简单的缓存封装 **collection.Cache**，简单使用方式如下

```go
// 初始化 cache，其中 WithLimit 可以指定最大缓存的数量
c, err := collection.NewCache(time.Minute, collection.WithLimit(10000))
if err != nil {
  panic(err)
}

// 设置缓存
c.Set("key", user)

// 获取缓存，ok：是否存在
v, ok := c.Get("key")

// 删除缓存
c.Del("key")

// 获取缓存，如果 key 不存在的，则会调用 func 去生成缓存
v, err := c.Take("key", func() (interface{}, error) {
  return user, nil
})
```

cache 实现的建的功能包括

* 缓存自动失效，可以指定过期时间
* 缓存大小限制，可以指定缓存个数
* 缓存增删改
* 缓存命中率统计
* 并发安全
* 缓存击穿

实现原理：
Cache 自动失效，是采用 TimingWheel(https://github.com/tal-tech/go-zero/blob/master/core/collection/timingwheel.go) 进行管理的

``` go
timingWheel, err := NewTimingWheel(time.Second, slots, func(k, v interface{}) {
		key, ok := k.(string)
		if !ok {
			return
		}

		cache.Del(key)
})
```

Cache 大小限制，是采用 LRU 淘汰策略，在新增缓存的时候会去检查是否已经超出过限制，具体代码在 keyLru 中实现

``` go
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
```

Cache 的命中率统计，是在代码中实现 cacheStat,在缓存命中丢失的时候自动统计，并且会定时打印使用的命中率, qps 等状态.

打印的具体效果如下

```go
cache(proc) - qpm: 2, hit_ratio: 50.0%, elements: 0, hit: 1, miss: 1
```

缓存击穿包含是使用 syncx.SharedCalls(https://github.com/tal-tech/go-zero/blob/master/core/syncx/sharedcalls.go) 进行实现的，就是将同时请求同一个 key 的请求, 关于 sharedcalls 后续会继续补充。 相关具体实现是在:

```go
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
```

本文主要介绍了go-zero框架中的 Cache 工具，在实际的项目中非常实用。用好工具对于提升服务性能和开发效率都有很大的帮助，希望本篇文章能给大家带来一些收获。
