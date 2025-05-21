package mathx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
)

func TestMaxInt(t *testing.T) {
	cases := []struct {
		a      int
		b      int
		expect int
	}{
		{
			a:      0,
			b:      1,
			expect: 1,
		},
		{
			a:      0,
			b:      -1,
			expect: 0,
		},
		{
			a:      1,
			b:      1,
			expect: 1,
		},
	}

	for _, each := range cases {
		each := each
		t.Run(stringx.Rand(), func(t *testing.T) {
			actual := MaxInt(each.a, each.b)
			assert.Equal(t, each.expect, actual)
		})
	}
}

func TestMinInt(t *testing.T) {
	cases := []struct {
		a      int
		b      int
		expect int
	}{
		{
			a:      0,
			b:      1,
			expect: 0,
		},
		{
			a:      0,
			b:      -1,
			expect: -1,
		},
		{
			a:      1,
			b:      1,
			expect: 1,
		},
	}

	for _, each := range cases {
		t.Run(stringx.Rand(), func(t *testing.T) {
			actual := MinInt(each.a, each.b)
			assert.Equal(t, each.expect, actual)
		})
	}
}
