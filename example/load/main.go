package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/core/executors"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/syncx"
	"gopkg.in/cheggaaa/pb.v1"
)

const (
	beta     = 0.9
	total    = 400
	interval = time.Second
	factor   = 5
)

var (
	seconds             = flag.Int("d", 400, "duration to go")
	flying              uint64
	avgFlyingAggressive float64
	aggressiveLock      syncx.SpinLock
	avgFlyingLazy       float64
	lazyLock            syncx.SpinLock
	avgFlyingBoth       float64
	bothLock            syncx.SpinLock
	lessWriter          *executors.LessExecutor
	passCounter         = collection.NewRollingWindow(50, time.Millisecond*100)
	rtCounter           = collection.NewRollingWindow(50, time.Millisecond*100)
	index               int32
)

func main() {
	flag.Parse()

	// only log 100 records
	lessWriter = executors.NewLessExecutor(interval * total / 100)

	fp, err := os.Create("result.csv")
	logx.Must(err)
	defer fp.Close()
	fmt.Fprintln(fp, "second,maxFlight,flying,agressiveAvgFlying,lazyAvgFlying,bothAvgFlying")

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	bar := pb.New(*seconds * 2).Start()
	var waitGroup sync.WaitGroup
	batchRequests := func(i int) {
		<-ticker.C
		requests := (i + 1) * factor
		func() {
			it := time.NewTicker(interval / time.Duration(requests))
			defer it.Stop()
			for j := 0; j < requests; j++ {
				<-it.C
				waitGroup.Add(1)
				go func() {
					issueRequest(fp, atomic.AddInt32(&index, 1))
					waitGroup.Done()
				}()
			}
			bar.Increment()
		}()
	}
	for i := 0; i < *seconds; i++ {
		batchRequests(i)
	}
	for i := *seconds; i > 0; i-- {
		batchRequests(i)
	}
	bar.Finish()
	waitGroup.Wait()
}

func issueRequest(writer io.Writer, idx int32) {
	v := atomic.AddUint64(&flying, 1)
	aggressiveLock.Lock()
	af := avgFlyingAggressive*beta + float64(v)*(1-beta)
	avgFlyingAggressive = af
	aggressiveLock.Unlock()
	bothLock.Lock()
	bf := avgFlyingBoth*beta + float64(v)*(1-beta)
	avgFlyingBoth = bf
	bothLock.Unlock()
	duration := time.Millisecond * time.Duration(rand.Int63n(10)+1)
	job(duration)
	passCounter.Add(1)
	rtCounter.Add(float64(duration) / float64(time.Millisecond))
	v1 := atomic.AddUint64(&flying, ^uint64(0))
	lazyLock.Lock()
	lf := avgFlyingLazy*beta + float64(v1)*(1-beta)
	avgFlyingLazy = lf
	lazyLock.Unlock()
	bothLock.Lock()
	bf = avgFlyingBoth*beta + float64(v1)*(1-beta)
	avgFlyingBoth = bf
	bothLock.Unlock()
	lessWriter.DoOrDiscard(func() {
		fmt.Fprintf(writer, "%d,%d,%d,%.2f,%.2f,%.2f\n", idx, maxFlight(), v, af, lf, bf)
	})
}

func job(duration time.Duration) {
	time.Sleep(duration)
}

func maxFlight() int64 {
	return int64(math.Max(1, float64(maxPass()*10)*(minRt()/1e3)))
}

func maxPass() int64 {
	var result float64 = 1

	passCounter.Reduce(func(b *collection.Bucket) {
		if b.Sum > result {
			result = b.Sum
		}
	})

	return int64(result)
}

func minRt() float64 {
	var result float64 = 1000

	rtCounter.Reduce(func(b *collection.Bucket) {
		if b.Count <= 0 {
			return
		}

		avg := math.Round(b.Sum / float64(b.Count))
		if avg < result {
			result = avg
		}
	})

	return result
}
