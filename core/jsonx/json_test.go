package jsonx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	var v = struct {
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
	var v = struct {
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
