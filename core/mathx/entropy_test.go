package mathx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcEntropy(t *testing.T) {
	const total = 1000
	const count = 100
	m := make(map[any]int, total)
	for i := 0; i < total; i++ {
		m[i] = count
	}
	assert.True(t, CalcEntropy(m) > .99)
}

func TestCalcEmptyEntropy(t *testing.T) {
	m := make(map[any]int)
	assert.Equal(t, float64(1), CalcEntropy(m))
}

func TestCalcDiffEntropy(t *testing.T) {
	const total = 1000
	m := make(map[any]int, total)
	for i := 0; i < total; i++ {
		m[i] = i
	}
	assert.True(t, CalcEntropy(m) < .99)
}
