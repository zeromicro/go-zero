# 防止缓存击穿之进程内共享调用

go-zero微服务框架中提供了许多开箱即用的工具，好的工具不仅能提升服务的性能而且还能提升代码的鲁棒性避免出错，实现代码风格的统一方便他人阅读等等。

本文主要讲述进程内共享调用神器[SharedCalls](https://github.com/tal-tech/go-zero/blob/master/core/syncx/sharedcalls.go)。  

## 使用场景

并发场景下，可能会有多个线程（协程）同时请求同一份资源，如果每个请求都要走一遍资源的请求过程，除了比较低效之外，还会对资源服务造成并发的压力。举一个具体例子，比如缓存失效，多个请求同时到达某服务请求某资源，该资源在缓存中已经失效，此时这些请求会继续访问DB做查询，会引起数据库压力瞬间增大。而使用SharedCalls可以使得同时多个请求只需要发起一次拿结果的调用，其他请求"坐享其成"，这种设计有效减少了资源服务的并发压力，可以有效防止缓存击穿。

高并发场景下，当某个热点key缓存失效后，多个请求会同时从数据库加载该资源，并保存到缓存，如果不做防范，可能会导致数据库被直接打死。针对这种场景，go-zero框架中已经提供了实现，具体可参看[sqlc](https://github.com/tal-tech/go-zero/blob/master/core/stores/sqlc/cachedsql.go)和[mongoc](https://github.com/tal-tech/go-zero/blob/master/core/stores/mongoc/cachedcollection.go)等实现代码。

为了简化演示代码，我们通过多个线程同时去获取一个id来模拟缓存的场景。如下：

```go
func main() {
  const round = 5
  var wg sync.WaitGroup
  barrier := syncx.NewSharedCalls()

  wg.Add(round)
  for i := 0; i < round; i++ {
    // 多个线程同时执行
    go func() {
      defer wg.Done()
      // 可以看到，多个线程在同一个key上去请求资源，获取资源的实际函数只会被调用一次
      val, err := barrier.Do("once", func() (interface{}, error) {
        // sleep 1秒，为了让多个线程同时取once这个key上的数据
        time.Sleep(time.Second)
        // 生成了一个随机的id
        return stringx.RandId(), nil
      })
      if err != nil {
        fmt.Println(err)
      } else {
        fmt.Println(val)
      }
    }()
  }

  wg.Wait()
}
```

运行，打印结果为：

```
837c577b1008a0db
837c577b1008a0db
837c577b1008a0db
837c577b1008a0db
837c577b1008a0db
```

可以看出，只要是同一个key上的同时发起的请求，都会共享同一个结果，对获取DB数据进缓存等场景特别有用，可以有效防止缓存击穿。

## 关键源码分析

- SharedCalls interface提供了Do和DoEx两种方法的抽象

  ```go
  // SharedCalls接口提供了Do和DoEx两种方法
  type SharedCalls interface {
    Do(key string, fn func() (interface{}, error)) (interface{}, error)
    DoEx(key string, fn func() (interface{}, error)) (interface{}, bool, error)
  }
  ```

- SharedCalls interface的具体实现sharedGroup

  ```go
  // call代表对指定资源的一次请求
  type call struct {
    wg  sync.WaitGroup  // 用于协调各个请求goroutine之间的资源共享
    val interface{}     // 用于保存请求的返回值
    err error           // 用于保存请求过程中发生的错误
  }
  
  type sharedGroup struct {
    calls map[string]*call
    lock  sync.Mutex
  }
  ```

- sharedGroup的Do方法

  - key参数：可以理解为资源的唯一标识。
  - fn参数：真正获取资源的方法。
  - 处理过程分析：

  ```go
  // 当多个请求同时使用Do方法请求资源时
  func (g *sharedGroup) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
    // 先申请加锁
    g.lock.Lock()
  
    // 根据key，获取对应的call结果,并用变量c保存
    if c, ok := g.calls[key]; ok {
      // 拿到call以后，释放锁，此处call可能还没有实际数据，只是一个空的内存占位
      g.lock.Unlock()
      // 调用wg.Wait，判断是否有其他goroutine正在申请资源，如果阻塞，说明有其他goroutine正在获取资源
      c.wg.Wait()
      // 当wg.Wait不再阻塞，表示资源获取已经结束，可以直接返回结果
      return c.val, c.err
    }

    // 没有拿到结果，则调用makeCall方法去获取资源，注意此处仍然是锁住的，可以保证只有一个goroutine可以调用makecall
    c := g.makeCall(key, fn)
    // 返回调用结果
    return c.val, c.err
  }
  ```
  
- sharedGroup的DoEx方法

  - 和Do方法类似，只是返回值中增加了布尔值表示值是调用makeCall方法直接获取的，还是取的共享成果

  ```go
  func (g *sharedGroup) DoEx(key string, fn func() (interface{}, error)) (val interface{}, fresh bool, err error) {
    g.lock.Lock()
    if c, ok := g.calls[key]; ok {
      g.lock.Unlock()
      c.wg.Wait()
      return c.val, false, c.err
    }

    c := g.makeCall(key, fn)
    return c.val, true, c.err
  }
  ```

- sharedGroup的makeCall方法

  - 该方法由Do和DoEx方法调用，是真正发起资源请求的方法。
  
  ```go
  // 进入makeCall的一定只有一个goroutine，因为要拿锁锁住的
  func (g *sharedGroup) makeCall(key string, fn func() (interface{}, error)) *call {
    // 创建call结构，用于保存本次请求的结果
    c := new(call)
    // wg加1，用于通知其他请求资源的goroutine等待本次资源获取的结束
    c.wg.Add(1)
    // 将用于保存结果的call放入map中，以供其他goroutine获取
    g.calls[key] = c
    // 释放锁，这样其他请求的goroutine才能获取call的内存占位
    g.lock.Unlock()
  
    defer func() {
      // delete key first, done later. can't reverse the order, because if reverse,
      // another Do call might wg.Wait() without get notified with wg.Done()
      g.lock.Lock()
      delete(g.calls, key)
      g.lock.Unlock()

      // 调用wg.Done，通知其他goroutine可以返回结果，这样本批次所有请求完成结果的共享
      c.wg.Done()
    }()
  
    // 调用fn方法，将结果填入变量c中
    c.val, c.err = fn()
    return c
  }
  ```

## 最后

本文主要介绍了go-zero框架中的 SharedCalls工具，对其应用场景和关键代码做了简单的梳理，希望本篇文章能给大家带来一些收获。
