package main

import (
	"fmt"
	"testing"

	"github.com/tal-tech/go-zero/core/fx"
)

func TestFxSplit(t *testing.T) {
	fx.Just(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Split(4).ForEach(func(item interface{}) {
		vals := item.([]interface{})
		fmt.Println(len(vals))
	})
}

func BenchmarkFx(b *testing.B) {
	type Mixed struct {
		Name   string
		Age    int
		Gender int
	}
	for i := 0; i < b.N; i++ {
		var mx Mixed
		fx.Parallel(func() {
			mx.Name = "hello"
		}, func() {
			mx.Age = 20
		}, func() {
			mx.Gender = 1
		})
	}
}
