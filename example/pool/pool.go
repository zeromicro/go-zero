package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/core/syncx"
)

func main() {
	var count int32
	var consumed int32
	pool := syncx.NewPool(80, func() interface{} {
		fmt.Printf("+ %d\n", atomic.AddInt32(&count, 1))
		return 1
	}, func(interface{}) {
		fmt.Printf("- %d\n", atomic.AddInt32(&count, -1))
	}, syncx.WithMaxAge(time.Second))

	var waitGroup sync.WaitGroup
	quit := make(chan lang.PlaceholderType)
	waitGroup.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer func() {
				waitGroup.Done()
				fmt.Println("routine quit")
			}()

			for {
				select {
				case <-quit:
					return
				default:
					x := pool.Get().(int)
					atomic.AddInt32(&consumed, 1)
					pool.Put(x)
				}
			}
		}()
	}

	bufio.NewReader(os.Stdin).ReadLine()
	close(quit)
	fmt.Println("quitted")
	waitGroup.Wait()
	fmt.Printf("consumed %d\n", atomic.LoadInt32(&consumed))
}
