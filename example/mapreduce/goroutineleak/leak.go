package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/mr"
	"github.com/tal-tech/go-zero/core/proc"
)

func dumpGoroutines() {
	dumpFile := "goroutines.dump"
	logx.Infof("Got dump goroutine signal, printing goroutine profile to %s", dumpFile)

	if f, err := os.Create(dumpFile); err != nil {
		logx.Errorf("Failed to dump goroutine profile, error: %v", err)
	} else {
		defer f.Close()
		pprof.Lookup("goroutine").WriteTo(f, 2)
	}
}

func main() {
	profiler := proc.StartProfile()
	defer profiler.Stop()

	done := make(chan lang.PlaceholderType)
	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println(runtime.NumGoroutine())
		}
	}()
	go func() {
		time.Sleep(time.Minute)
		dumpGoroutines()
		close(done)
	}()
	for {
		select {
		case <-done:
			return
		default:
			mr.MapReduce(func(source chan<- interface{}) {
				for i := 0; i < 100; i++ {
					source <- i
				}
			}, func(item interface{}, writer mr.Writer, cancel func(error)) {
				if item.(int) == 40 {
					cancel(errors.New("any"))
					return
				}
				writer.Write(item)
			}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
				list := make([]int, 0)
				for p := range pipe {
					list = append(list, p.(int))
				}
				writer.Write(list)
			})
		}
	}
}
