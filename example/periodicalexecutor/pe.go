package main

import (
	"time"

	"zero/core/executors"
)

func main() {
	exeutor := executors.NewBulkExecutor(func(items []interface{}) {
		println(len(items))
	}, executors.WithBulkTasks(10))
	for {
		exeutor.Add(1)
		time.Sleep(time.Millisecond * 90)
	}
}
