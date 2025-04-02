package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type data struct {
	input    string
	expected string
}

func TestTitle(t *testing.T) {
	list := []*data{
		{input: "_", expected: "_"},
		{input: "abc", expected: "Abc"},
		{input: "ABC", expected: "ABC"},
		{input: "", expected: ""},
		{input: " abc", expected: " abc"},
	}
	for _, e := range list {
		assert.Equal(t, e.expected, Title(e.input))
	}
}

func TestUntitle(t *testing.T) {
	list := []*data{
		{input: "_", expected: "_"},
		{input: "Abc", expected: "abc"},
		{input: "ABC", expected: "aBC"},
		{input: "", expected: ""},
		{input: " abc", expected: " abc"},
	}

	for _, e := range list {
		assert.Equal(t, e.expected, Untitle(e.input))
	}
}

func TestIndex(t *testing.T) {
	list := []string{"a", "b", "c"}
	assert.Equal(t, 1, Index(list, "b"))
	assert.Equal(t, -1, Index(list, "d"))
}

func TestSafeString(t *testing.T) {
	list := []*data{
		{input: "_", expected: "_"},
		{input: "a-b-c", expected: "a_b_c"},
		{input: "123abc", expected: "_123abc"},
		{input: "汉abc", expected: "_abc"},
		{input: "汉a字", expected: "_a_"},
		{input: "キャラクターabc", expected: "______abc"},
		{input: "-a_B-C", expected: "_a_B_C"},
		{input: "a_B C", expected: "a_B_C"},
		{input: "A#B#C", expected: "A_B_C"},
		{input: "_123", expected: "_123"},
		{input: "", expected: ""},
		{input: "\t", expected: "_"},
		{input: "\n", expected: "_"},
	}
	for _, e := range list {
		assert.Equal(t, e.expected, SafeString(e.input))
	}
}

func TestEscapeGoKeyword(t *testing.T) {
	for k := range goKeyword {
		assert.Equal(t, goKeyword[k], EscapeGolangKeyword(k))
		assert.False(t, isGolangKeyword(strings.Title(k)))
	}
}

func TestFieldsAndTrimSpace(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{input: " a b c ", expected: []string{"a", "b", "c"}},
		{input: "  a   b   c  ", expected: []string{"a", "b", "c"}},
		{input: "\ta\tb\tc\t", expected: []string{"a", "b", "c"}},
		{input: "\na\nb\nc\n", expected: []string{"a", "b", "c"}},
		{input: "a b c", expected: []string{"a", "b", "c"}},
		{input: "a\tb\tc", expected: []string{"a", "b", "c"}},
		{input: "a\nb\nc", expected: []string{"a", "b", "c"}},
		{input: "", expected: []string{}},
		{input: "   ", expected: []string{}},
	}

	for _, tc := range testCases {
		result := FieldsAndTrimSpace(tc.input, func(r rune) bool {
			return r == ' ' || r == '\t' || r == '\n'
		})
		assert.Equal(t, tc.expected, result)
	}
}

func TestUnquote(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: `"hello"`, expected: `hello`},
		{input: "`world`", expected: `world`},
		{input: `'test'`, expected: `test`},
		{input: `"foo'bar"`, expected: `foobar`},
		{input: "", expected: ""},
	}

	for _, tc := range testCases {
		result := Unquote(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}
