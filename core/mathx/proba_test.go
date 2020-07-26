package mathx

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrueOnProba(t *testing.T) {
	const proba = math.Pi / 10
	const total = 100000
	const epsilon = 0.05
	var count int
	p := NewProba()
	for i := 0; i < total; i++ {
		if p.TrueOnProba(proba) {
			count++
		}
	}

	ratio := float64(count) / float64(total)
	assert.InEpsilon(t, proba, ratio, epsilon)
}
