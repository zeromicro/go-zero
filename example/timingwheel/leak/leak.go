package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/tal-tech/go-zero/core/cmdline"
	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/core/proc"
)

const numItems = 1000000

var round = flag.Int("r", 3, "rounds to go")

func main() {
	defer proc.StartProfile().Stop()

	flag.Parse()

	fmt.Println(getMemUsage())
	for i := 0; i < *round; i++ {
		do()
	}
	cmdline.EnterToContinue()
}

func do() {
	tw, err := collection.NewTimingWheel(time.Second, 100, execute)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < numItems; i++ {
		key := strconv.Itoa(i)
		tw.SetTimer(key, key, time.Second*5)
	}

	fmt.Println(getMemUsage())
}

func execute(k, v interface{}) {
}

func getMemUsage() string {
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For more info, see: https://golang.org/pkg/runtime/#MemStats
	return fmt.Sprintf("Heap Alloc = %dMiB", toMiB(m.HeapAlloc))
}

func toMiB(b uint64) uint64 {
	return b / 1024 / 1024
}
