package main

import (
	"fmt"

	"github.com/tal-tech/go-zero/core/fx"
)

func main() {
	result, err := fx.From(func(source chan<- interface{}) {
		for i := 0; i < 10; i++ {
			source <- i
			source <- i
		}
	}).Map(func(item interface{}) interface{} {
		i := item.(int)
		return i * i
	}).Filter(func(item interface{}) bool {
		i := item.(int)
		return i%2 == 0
	}).Distinct(func(item interface{}) interface{} {
		return item
	}).Reduce(func(pipe <-chan interface{}) (interface{}, error) {
		var result int
		for item := range pipe {
			i := item.(int)
			result += i
		}
		return result, nil
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}
