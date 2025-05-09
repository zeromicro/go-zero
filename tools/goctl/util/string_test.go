package util

import (
	"strings"
	"testing"
	"unicode"

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
		name     string
		input    string
		delimiter func(r rune) bool
		expected []string
	}{
		{
			name:     "Comma-separated values",
			input:    "a, b, c",
			delimiter: func(r rune) bool { return r == ',' },
			expected: []string{"a", " b", " c"},
		},
		{
			name:     "Space-separated values",
			input:    "a b c",
			delimiter: unicode.IsSpace,
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Mixed whitespace",
			input:    "a\tb\nc",
			delimiter: unicode.IsSpace,
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Empty input",
			input:    "",
			delimiter: unicode.IsSpace,
			expected: []string(nil),
		},
		{
			name:     "Trailing and leading spaces",
			input:    "  a , b , c  ",
			delimiter: func(r rune) bool { return r == ',' },
			expected: []string{"  a ", " b ", " c  "},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FieldsAndTrimSpace(tc.input, tc.delimiter)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUnquote(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: `"hello"`, expected: `hello`},
		{input: "`world`", expected: `world`},
		{input: `"foo'bar"`, expected: `foo'bar`},
		{input: "", expected: ""},
	}

	for _, tc := range testCases {
		result := Unquote(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}
