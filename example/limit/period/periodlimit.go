package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/limit"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

const seconds = 5

var (
	rdx     = flag.String("redis", "localhost:6379", "the redis, default localhost:6379")
	rdxType = flag.String("redisType", "node", "the redis type, default node")
	rdxPass = flag.String("redisPass", "", "the redis password")
	rdxKey  = flag.String("redisKey", "rate", "the redis key, default rate")
	threads = flag.Int("threads", runtime.NumCPU(), "the concurrent threads, default to cores")
)

func main() {
	flag.Parse()

	store := redis.NewRedis(*rdx, *rdxType, *rdxPass)
	fmt.Println(store.Ping())
	lmt := limit.NewPeriodLimit(seconds, 5, store, *rdxKey)
	timer := time.NewTimer(time.Second * seconds)
	quit := make(chan struct{})
	defer timer.Stop()
	go func() {
		<-timer.C
		close(quit)
	}()

	var allowed, denied int32
	var wait sync.WaitGroup
	for i := 0; i < *threads; i++ {
		wait.Add(1)
		go func() {
			for {
				select {
				case <-quit:
					wait.Done()
					return
				default:
					if v, err := lmt.Take(strconv.FormatInt(int64(i), 10)); err == nil && v == limit.Allowed {
						atomic.AddInt32(&allowed, 1)
					} else if err != nil {
						log.Fatal(err)
					} else {
						atomic.AddInt32(&denied, 1)
					}
				}
			}
		}()
	}

	wait.Wait()
	fmt.Printf("allowed: %d, denied: %d, qps: %d\n", allowed, denied, (allowed+denied)/seconds)
}
