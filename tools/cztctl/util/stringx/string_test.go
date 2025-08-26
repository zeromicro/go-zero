package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_IsEmptyOrSpace(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{
			want: true,
		},
		{
			input: " ",
			want:  true,
		},
		{
			input: "\t",
			want:  true,
		},
		{
			input: "\n",
			want:  true,
		},
		{
			input: "\f",
			want:  true,
		},
		{
			input: "		",
			want:  true,
		},
	}
	for _, v := range cases {
		s := From(v.input)
		assert.Equal(t, v.want, s.IsEmptyOrSpace())
	}
}

func TestString_Snake2Camel(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "__",
			want:  "",
		},
		{
			input: "go_zero",
			want:  "GoZero",
		},
		{
			input: "の_go_zero",
			want:  "のGoZero",
		},
		{
			input: "goZero",
			want:  "GoZero",
		},
		{
			input: "goZero",
			want:  "GoZero",
		},
		{
			input: "goZero_",
			want:  "GoZero",
		},
		{
			input: "go_Zero_",
			want:  "GoZero",
		},
		{
			input: "_go_Zero_",
			want:  "GoZero",
		},
	}
	for _, c := range cases {
		ret := From(c.input).ToCamel()
		assert.Equal(t, c.want, ret)
	}
}

func TestString_Camel2Snake(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "goZero",
			want:  "go_zero",
		},
		{
			input: "Gozero",
			want:  "gozero",
		},
		{
			input: "GoZero",
			want:  "go_zero",
		},
		{
			input: "Go_Zero",
			want:  "go__zero",
		},
	}
	for _, c := range cases {
		ret := From(c.input).ToSnake()
		assert.Equal(t, c.want, ret)
	}
}

func TestTitle(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "go zero",
			want:  "Go Zero",
		},
		{
			input: "goZero",
			want:  "GoZero",
		},
		{
			input: "GoZero",
			want:  "GoZero",
		},
		{
			input: "の go zero",
			want:  "の Go Zero",
		},
		{
			input: "Gozero",
			want:  "Gozero",
		},
		{
			input: "Go_zero",
			want:  "Go_zero",
		},
		{
			input: "go_zero",
			want:  "Go_zero",
		},
		{
			input: "Go_Zero",
			want:  "Go_Zero",
		},
	}
	for _, c := range cases {
		ret := From(c.input).Title()
		assert.Equal(t, c.want, ret)
	}
}

func TestUntitle(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			input: "go zero",
			want:  "go zero",
		},
		{
			input: "GoZero",
			want:  "goZero",
		},
		{
			input: "Gozero",
			want:  "gozero",
		},
		{
			input: "Go_zero",
			want:  "go_zero",
		},
		{
			input: "go_zero",
			want:  "go_zero",
		},
		{
			input: "Go_Zero",
			want:  "go_Zero",
		},
	}
	for _, c := range cases {
		ret := From(c.input).Untitle()
		assert.Equal(t, c.want, ret)
	}
}

func TestContainsAny(t *testing.T) {
	type args struct {
		s     string
		runes []rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "runes is empty",
			args: args{
				s:     "test",
				runes: []rune{},
			},
			want: true,
		},
		{
			name: "s is empty and runes is not empty",
			args: args{
				s:     "",
				runes: []rune{'a', 'b', 'c'},
			},
			want: false,
		},
		{
			name: "s contains runes",
			args: args{
				s:     "hello",
				runes: []rune{'e', 'f'},
			},
			want: true,
		},
		{
			name: "s does not contain runes",
			args: args{
				s:     "hello",
				runes: []rune{'x', 'y'},
			},
			want: false,
		},
		{
			name: "s and runes both have one matching character",
			args: args{
				s:     "a",
				runes: []rune{'a'},
			},
			want: true,
		},
		{
			name: "s and runes both have one non-matching character",
			args: args{
				s:     "a",
				runes: []rune{'b'},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ContainsAny(tt.args.s, tt.args.runes...), "ContainsAny(%v, %v)", tt.args.s, tt.args.runes)
		})
	}
}

func TestContainsWhiteSpace(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "contains space",
			args: args{s: "hello world"},
			want: true,
		},
		{
			name: "contains newline",
			args: args{s: "hello\nworld"},
			want: true,
		},
		{
			name: "contains tab",
			args: args{s: "hello\tworld"},
			want: true,
		},
		{
			name: "contains form feed",
			args: args{s: "hello\fworld"},
			want: true,
		},
		{
			name: "contains vertical tab",
			args: args{s: "hello\vworld"},
			want: true,
		},
		{
			name: "no whitespace",
			args: args{s: "helloworld"},
			want: false,
		},
		{
			name: "empty string",
			args: args{s: ""},
			want: false,
		},
		{
			name: "only whitespace",
			args: args{s: "  \t\n\f\v"},
			want: true,
		},
		{
			name: "contains non-standard whitespace",
			args: args{s: "hello\u00A0world"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ContainsWhiteSpace(tt.args.s), "ContainsWhiteSpace(%v)", tt.args.s)
		})
	}
}
