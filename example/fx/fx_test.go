package main

import (
	"testing"

	"github.com/tal-tech/go-zero/core/fx"
)

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
