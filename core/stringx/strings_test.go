package stringx

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotEmpty(t *testing.T) {
	cases := []struct {
		args   []string
		expect bool
	}{
		{
			args:   []string{"a", "b", "c"},
			expect: true,
		},
		{
			args:   []string{"a", "", "c"},
			expect: false,
		},
		{
			args:   []string{"a"},
			expect: true,
		},
		{
			args:   []string{""},
			expect: false,
		},
		{
			args:   []string{},
			expect: true,
		},
	}

	for _, each := range cases {
		t.Run(path.Join(each.args...), func(t *testing.T) {
			assert.Equal(t, each.expect, NotEmpty(each.args...))
		})
	}
}

func TestContainsString(t *testing.T) {
	cases := []struct {
		slice  []string
		value  string
		expect bool
	}{
		{[]string{"1"}, "1", true},
		{[]string{"1"}, "2", false},
		{[]string{"1", "2"}, "1", true},
		{[]string{"1", "2"}, "3", false},
		{nil, "3", false},
		{nil, "", false},
	}

	for _, each := range cases {
		t.Run(path.Join(each.slice...), func(t *testing.T) {
			actual := Contains(each.slice, each.value)
			assert.Equal(t, each.expect, actual)
		})
	}
}

func TestFilter(t *testing.T) {
	cases := []struct {
		input   string
		ignores []rune
		expect  string
	}{
		{``, nil, ``},
		{`abcd`, nil, `abcd`},
		{`ab,cd,ef`, []rune{','}, `abcdef`},
		{`ab, cd,ef`, []rune{',', ' '}, `abcdef`},
		{`ab, cd, ef`, []rune{',', ' '}, `abcdef`},
		{`ab, cd, ef, `, []rune{',', ' '}, `abcdef`},
	}

	for _, each := range cases {
		t.Run(each.input, func(t *testing.T) {
			actual := Filter(each.input, func(r rune) bool {
				for _, x := range each.ignores {
					if x == r {
						return true
					}
				}
				return false
			})
			assert.Equal(t, each.expect, actual)
		})
	}
}

func TestRemove(t *testing.T) {
	cases := []struct {
		input  []string
		remove []string
		expect []string
	}{
		{
			input:  []string{"a", "b", "a", "c"},
			remove: []string{"a", "b"},
			expect: []string{"c"},
		},
		{
			input:  []string{"b", "c"},
			remove: []string{"a"},
			expect: []string{"b", "c"},
		},
		{
			input:  []string{"b", "a", "c"},
			remove: []string{"a"},
			expect: []string{"b", "c"},
		},
		{
			input:  []string{},
			remove: []string{"a"},
			expect: []string{},
		},
	}

	for _, each := range cases {
		t.Run(path.Join(each.input...), func(t *testing.T) {
			assert.ElementsMatch(t, each.expect, Remove(each.input, each.remove...))
		})
	}
}

func TestReverse(t *testing.T) {
	cases := []struct {
		input  string
		expect string
	}{
		{
			input:  "abcd",
			expect: "dcba",
		},
		{
			input:  "",
			expect: "",
		},
		{
			input:  "我爱中国",
			expect: "国中爱我",
		},
	}

	for _, each := range cases {
		t.Run(each.input, func(t *testing.T) {
			assert.Equal(t, each.expect, Reverse(each.input))
		})
	}
}

func TestSubstr(t *testing.T) {
	cases := []struct {
		input  string
		start  int
		stop   int
		err    error
		expect string
	}{
		{
			input:  "abcdefg",
			start:  1,
			stop:   4,
			expect: "bcd",
		},
		{
			input:  "我爱中国3000遍，even more",
			start:  1,
			stop:   9,
			expect: "爱中国3000遍",
		},
		{
			input:  "abcdefg",
			start:  -1,
			stop:   4,
			err:    ErrInvalidStartPosition,
			expect: "",
		},
		{
			input:  "abcdefg",
			start:  100,
			stop:   4,
			err:    ErrInvalidStartPosition,
			expect: "",
		},
		{
			input:  "abcdefg",
			start:  1,
			stop:   -1,
			err:    ErrInvalidStopPosition,
			expect: "",
		},
		{
			input:  "abcdefg",
			start:  1,
			stop:   100,
			err:    ErrInvalidStopPosition,
			expect: "",
		},
	}

	for _, each := range cases {
		t.Run(each.input, func(t *testing.T) {
			val, err := Substr(each.input, each.start, each.stop)
			assert.Equal(t, each.err, err)
			if err == nil {
				assert.Equal(t, each.expect, val)
			}
		})
	}
}

func TestTakeOne(t *testing.T) {
	cases := []struct {
		valid  string
		or     string
		expect string
	}{
		{"", "", ""},
		{"", "1", "1"},
		{"1", "", "1"},
		{"1", "2", "1"},
	}

	for _, each := range cases {
		t.Run(each.valid, func(t *testing.T) {
			actual := TakeOne(each.valid, each.or)
			assert.Equal(t, each.expect, actual)
		})
	}
}

func TestTakeWithPriority(t *testing.T) {
	tests := []struct {
		fns    []func() string
		expect string
	}{
		{
			fns: []func() string{
				func() string {
					return "first"
				},
				func() string {
					return "second"
				},
				func() string {
					return "third"
				},
			},
			expect: "first",
		},
		{
			fns: []func() string{
				func() string {
					return ""
				},
				func() string {
					return "second"
				},
				func() string {
					return "third"
				},
			},
			expect: "second",
		},
		{
			fns: []func() string{
				func() string {
					return ""
				},
				func() string {
					return ""
				},
				func() string {
					return "third"
				},
			},
			expect: "third",
		},
		{
			fns: []func() string{
				func() string {
					return ""
				},
				func() string {
					return ""
				},
				func() string {
					return ""
				},
			},
			expect: "",
		},
	}

	for _, test := range tests {
		t.Run(RandId(), func(t *testing.T) {
			val := TakeWithPriority(test.fns...)
			assert.Equal(t, test.expect, val)
		})
	}
}

func TestUnion(t *testing.T) {
	first := []string{
		"one",
		"two",
		"three",
	}
	second := []string{
		"zero",
		"two",
		"three",
		"four",
	}
	union := Union(first, second)
	contains := func(v string) bool {
		for _, each := range union {
			if v == each {
				return true
			}
		}

		return false
	}
	assert.Equal(t, 5, len(union))
	assert.True(t, contains("zero"))
	assert.True(t, contains("one"))
	assert.True(t, contains("two"))
	assert.True(t, contains("three"))
	assert.True(t, contains("four"))
}
