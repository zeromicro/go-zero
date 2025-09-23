package jsonx

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	v := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}
	bs, err := Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"John","age":30}`, string(bs))
}

func TestMarshalToString(t *testing.T) {
	v := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}
	toString, err := MarshalToString(v)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"John","age":30}`, toString)

	_, err = MarshalToString(make(chan int))
	assert.NotNil(t, err)
}

func TestUnmarshal(t *testing.T) {
	const s = `{"name":"John","age":30}`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := Unmarshal([]byte(s), &v)
	assert.Nil(t, err)
	assert.Equal(t, "John", v.Name)
	assert.Equal(t, 30, v.Age)
}

func TestUnmarshalError(t *testing.T) {
	const s = `{"name":"John","age":30`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := Unmarshal([]byte(s), &v)
	assert.NotNil(t, err)
}

func TestUnmarshalFromString(t *testing.T) {
	const s = `{"name":"John","age":30}`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := UnmarshalFromString(s, &v)
	assert.Nil(t, err)
	assert.Equal(t, "John", v.Name)
	assert.Equal(t, 30, v.Age)
}

func TestUnmarshalFromStringError(t *testing.T) {
	const s = `{"name":"John","age":30`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := UnmarshalFromString(s, &v)
	assert.NotNil(t, err)
}

func TestUnmarshalFromRead(t *testing.T) {
	const s = `{"name":"John","age":30}`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := UnmarshalFromReader(strings.NewReader(s), &v)
	assert.Nil(t, err)
	assert.Equal(t, "John", v.Name)
	assert.Equal(t, 30, v.Age)
}

func TestUnmarshalFromReaderError(t *testing.T) {
	const s = `{"name":"John","age":30`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := UnmarshalFromReader(strings.NewReader(s), &v)
	assert.NotNil(t, err)
}

func Test_doMarshalJson(t *testing.T) {
	type args struct {
		v any
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "nil",
			args:    args{nil},
			want:    []byte("null"),
			wantErr: assert.NoError,
		},
		{
			name:    "string",
			args:    args{"hello"},
			want:    []byte(`"hello"`),
			wantErr: assert.NoError,
		},
		{
			name:    "int",
			args:    args{42},
			want:    []byte("42"),
			wantErr: assert.NoError,
		},
		{
			name:    "bool",
			args:    args{true},
			want:    []byte("true"),
			wantErr: assert.NoError,
		},
		{
			name: "struct",
			args: args{
				struct {
					Name string `json:"name"`
				}{Name: "test"},
			},
			want:    []byte(`{"name":"test"}`),
			wantErr: assert.NoError,
		},
		{
			name:    "slice",
			args:    args{[]int{1, 2, 3}},
			want:    []byte("[1,2,3]"),
			wantErr: assert.NoError,
		},
		{
			name:    "map",
			args:    args{map[string]int{"a": 1, "b": 2}},
			want:    []byte(`{"a":1,"b":2}`),
			wantErr: assert.NoError,
		},
		{
			name:    "unmarshalable type",
			args:    args{complex(1, 2)},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "channel type",
			args:    args{make(chan int)},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name:    "url with query params",
			args:    args{"https://example.com/api?name=test&age=25"},
			want:    []byte(`"https://example.com/api?name=test&age=25"`),
			wantErr: assert.NoError,
		},
		{
			name:    "url with encoded query params",
			args:    args{"https://example.com/api?data=hello%20world&special=%26%3D"},
			want:    []byte(`"https://example.com/api?data=hello%20world&special=%26%3D"`),
			wantErr: assert.NoError,
		},
		{
			name:    "url with multiple query params",
			args:    args{"http://localhost:8080/users?page=1&limit=10&sort=name&order=asc"},
			want:    []byte(`"http://localhost:8080/users?page=1&limit=10&sort=name&order=asc"`),
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.v)
			if !tt.wantErr(t, err, fmt.Sprintf("Marshal(%v)", tt.args.v)) {
				return
			}

			assert.Equalf(t, string(tt.want), string(got), "Marshal(%v)", tt.args.v)
		})
	}
}

func TestMarshalWithBuffer(t *testing.T) {
	v := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}

	buf := &bytes.Buffer{}
	bs, err := MarshalWithBuffer(v, buf)
	assert.Nil(t, err)
	assert.Equal(t, `{"name":"John","age":30}`, string(bs))

	// Test that buffer content is as expected (includes trailing newline)
	bufContent := buf.String()
	assert.Equal(t, "{\"name\":\"John\",\"age\":30}\n", bufContent)
}

func TestMarshalWithBufferError(t *testing.T) {
	buf := &bytes.Buffer{}
	_, err := MarshalWithBuffer(make(chan int), buf)
	assert.NotNil(t, err)
}

func TestMarshalWithBufferReuse(t *testing.T) {
	buf := &bytes.Buffer{}

	// First marshal
	v1 := map[string]string{"key": "value1"}
	bs1, err := MarshalWithBuffer(v1, buf)
	assert.Nil(t, err)
	assert.Equal(t, `{"key":"value1"}`, string(bs1))

	// Reset buffer and reuse
	buf.Reset()
	v2 := map[string]string{"key": "value2"}
	bs2, err := MarshalWithBuffer(v2, buf)
	assert.Nil(t, err)
	assert.Equal(t, `{"key":"value2"}`, string(bs2))
}

func TestMarshalWithBufferMultipleTypes(t *testing.T) {
	tests := []struct {
		name string
		args interface{}
		want string
	}{
		{
			name: "nil",
			args: nil,
			want: "null",
		},
		{
			name: "string",
			args: "hello",
			want: `"hello"`,
		},
		{
			name: "int",
			args: 42,
			want: "42",
		},
		{
			name: "bool",
			args: true,
			want: "true",
		},
		{
			name: "struct",
			args: struct {
				Name string `json:"name"`
			}{Name: "test"},
			want: `{"name":"test"}`,
		},
		{
			name: "slice",
			args: []int{1, 2, 3},
			want: "[1,2,3]",
		},
		{
			name: "map",
			args: map[string]int{"a": 1, "b": 2},
			want: `{"a":1,"b":2}`,
		},
		{
			name: "url with special characters",
			args: "https://example.com/api?name=test&age=25",
			want: `"https://example.com/api?name=test&age=25"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			got, err := MarshalWithBuffer(tt.args, buf)
			assert.Nil(t, err, "MarshalWithBuffer should not return error for %v", tt.args)
			assert.Equal(t, tt.want, string(got), "MarshalWithBuffer(%v)", tt.args)
		})
	}
}
