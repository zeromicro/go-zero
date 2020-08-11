package main

import (
	"fmt"
	"time"

	"github.com/tal-tech/go-zero/core/executors"
)

func main() {
	exeutor := executors.NewBulkExecutor(func(items []interface{}) {
		fmt.Println(len(items))
	}, executors.WithBulkTasks(10))
	for {
		exeutor.Add(1)
		time.Sleep(time.Millisecond * 90)
	}
}
