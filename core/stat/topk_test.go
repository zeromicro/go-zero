package stat

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	numSamples = 10000
	topNum     = 100
)

var samples []Task

func init() {
	for i := 0; i < numSamples; i++ {
		task := Task{
			Duration: time.Duration(rand.Int63()),
		}
		samples = append(samples, task)
	}
}

func TestTopK(t *testing.T) {
	tasks := []Task{
		{false, 1, "a"},
		{false, 4, "a"},
		{false, 2, "a"},
		{false, 5, "a"},
		{false, 9, "a"},
		{false, 10, "a"},
		{false, 12, "a"},
		{false, 3, "a"},
		{false, 6, "a"},
		{false, 11, "a"},
		{false, 8, "a"},
	}

	result := topK(tasks, 3)
	if len(result) != 3 {
		t.Fail()
	}

	set := make(map[time.Duration]struct{})
	for _, each := range result {
		set[each.Duration] = struct{}{}
	}

	for _, v := range []time.Duration{10, 11, 12} {
		_, ok := set[v]
		assert.True(t, ok)
	}
}

func BenchmarkTopkHeap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		topK(samples, topNum)
	}
}
