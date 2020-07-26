package main

import "zero/core/threading"

func main() {
	q := threading.NewTaskRunner(5)
	q.Schedule(func() {
		panic("hello")
	})
	select {}
}
