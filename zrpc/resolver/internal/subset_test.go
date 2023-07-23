package internal

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/mathx"
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
		test := test
		t.Run(test.name, func(t *testing.T) {
			var vals []string
			for i := 0; i < test.set; i++ {
				vals = append(vals, strconv.Itoa(i))
			}

			m := make(map[any]int)
			for i := 0; i < 1000; i++ {
				set := subset(append([]string(nil), vals...), test.sub)
				if test.sub < test.set {
					assert.Equal(t, test.sub, len(set))
				} else {
					assert.Equal(t, test.set, len(set))
				}

				for _, val := range set {
					m[val]++
				}
			}

			assert.True(t, mathx.CalcEntropy(m) > 0.95)
		})
	}
}
