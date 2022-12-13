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
