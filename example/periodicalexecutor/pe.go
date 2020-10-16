package main

import (
	"fmt"
	"time"

	"github.com/tal-tech/go-zero/core/executors"
)

func main() {
	executor := executors.NewBulkExecutor(func(items []interface{}) {
		fmt.Println(len(items))
	}, executors.WithBulkTasks(10))
	for {
		if err := executor.Add(1); err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(time.Millisecond * 90)
	}
}
