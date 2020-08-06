package resolver

import (
	"strconv"
	"testing"

	"zero/core/mathx"

	"github.com/stretchr/testify/assert"
)

func TestSubset(t *testing.T) {
	tests := []struct {
		name string
		set  int
		sub  int
	}{
		{
			name: "more vals to subset",
			set:  100,
			sub:  36,
		},
		{
			name: "less vals to subset",
			set:  100,
			sub:  200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var vals []string
			for i := 0; i < test.set; i++ {
				vals = append(vals, strconv.Itoa(i))
			}

			m := make(map[interface{}]int)
			for i := 0; i < 1000; i++ {
				set := subset(append([]string(nil), vals...), test.sub)
				for _, val := range set {
					m[val]++
				}
			}

			assert.True(t, mathx.CalcEntropy(m) > 0.95)
		})
	}
}

func TestSubsetLess(t *testing.T) {
	var vals []string
	for i := 0; i < 100; i++ {
		vals = append(vals, strconv.Itoa(i))
	}

	m := make(map[interface{}]int)
	for i := 0; i < 1000; i++ {
		set := subset(append([]string(nil), vals...), 200)
		for _, val := range set {
			m[val]++
		}
	}

	assert.True(t, mathx.CalcEntropy(m) > 0.95)
}
