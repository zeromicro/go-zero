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
	v, err := Unmarshal[struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}]([]byte(s))
	assert.Nil(t, err)
	assert.Equal(t, "John", v.Name)
	assert.Equal(t, 30, v.Age)
}

func TestUnmarshalError(t *testing.T) {
	const s = `{"name":"John","age":30`
	_, err := Unmarshal[struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}]([]byte(s))
	assert.NotNil(t, err)
}

func TestUnmarshalFromString(t *testing.T) {
	const s = `{"name":"John","age":30}`
	v, err := UnmarshalFromString[struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}](s)
	assert.Nil(t, err)
	assert.Equal(t, "John", v.Name)
	assert.Equal(t, 30, v.Age)
}

func TestUnmarshalFromStringError(t *testing.T) {
	const s = `{"name":"John","age":30`
	_, err := UnmarshalFromString[struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}](s)
	assert.NotNil(t, err)
}

func TestUnmarshalFromRead(t *testing.T) {
	const s = `{"name":"John","age":30}`
	v, err := UnmarshalFromReader[struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}](strings.NewReader(s))
	assert.Nil(t, err)
	assert.Equal(t, "John", v.Name)
	assert.Equal(t, 30, v.Age)
}

func TestUnmarshalFromReaderError(t *testing.T) {
	const s = `{"name":"John","age":30`
	_, err := UnmarshalFromReader[struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}](strings.NewReader(s))
	assert.NotNil(t, err)
}
