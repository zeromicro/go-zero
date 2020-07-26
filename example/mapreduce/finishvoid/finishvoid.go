package main

import (
	"fmt"
	"time"

	"zero/core/mapreduce"
	"zero/core/timex"
)

func main() {
	start := timex.Now()

	mapreduce.FinishVoid(func() {
		time.Sleep(time.Second)
	}, func() {
		time.Sleep(time.Second * 5)
	}, func() {
		time.Sleep(time.Second * 10)
	}, func() {
		time.Sleep(time.Second * 6)
	}, func() {
		if err := mapreduce.Finish(func() error {
			time.Sleep(time.Second)
			return nil
		}, func() error {
			time.Sleep(time.Second * 10)
			return nil
		}); err != nil {
			fmt.Println(err)
		}
	})

	fmt.Println(timex.Since(start))
}
