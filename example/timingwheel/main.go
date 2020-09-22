package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/collection"
)

const interval = time.Minute

var traditional = flag.Bool("traditional", false, "enable traditional mode")

func main() {
	flag.Parse()

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())
			}
		}
	}()

	if *traditional {
		traditionalMode()
	} else {
		timingWheelMode()
	}
}

func timingWheelMode() {
	var count uint64
	tw, err := collection.NewTimingWheel(time.Second, 600, func(key, value interface{}) {
		job(&count)
	})
	if err != nil {
		log.Fatal(err)
	}

	defer tw.Stop()
	for i := 0; ; i++ {
		tw.SetTimer(i, i, interval)
		time.Sleep(time.Millisecond)
	}
}

func traditionalMode() {
	var count uint64
	for {
		go func() {
			timer := time.NewTimer(interval)
			defer timer.Stop()

			select {
			case <-timer.C:
				job(&count)
			}
		}()

		time.Sleep(time.Millisecond)
	}
}

func job(count *uint64) {
	v := atomic.AddUint64(count, 1)
	if v%1000 == 0 {
		fmt.Println(v)
	}
}
