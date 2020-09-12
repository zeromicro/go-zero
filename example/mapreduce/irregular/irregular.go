package main

import (
	"fmt"
	"time"

	"github.com/tal-tech/go-zero/core/mr"
)

func main() {
	mr.MapReduceVoid(func(source chan<- interface{}) {
		for i := 0; i < 10; i++ {
			source <- i
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		i := item.(int)
		if i == 0 {
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(5 * time.Second)
		}
		writer.Write(i)
	}, func(pipe <-chan interface{}, cancel func(error)) {
		for i := range pipe {
			fmt.Println(i)
		}
	})
}
