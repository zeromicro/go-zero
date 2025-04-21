package jsoncode

import (
	"reflect"
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

func TestSetMarshal(t *testing.T) {
	SetMarshalFn(func(v any) ([]byte, error) {
		return []byte{}, nil
	})

	var v = struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}
	bs, err := Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, ``, string(bs))
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

func TestSetUnmarshal(t *testing.T) {
	SetUnmarshalFn(func(by []byte, v any) error {
		vValue := reflect.ValueOf(v).Elem()
		vValue.FieldByName("Name").SetString("SetJohn")
		vValue.FieldByName("Age").SetInt(40)

		return nil
	})

	const s = `{"name":"John","age":30}`
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	err := Unmarshal([]byte(s), &v)
	assert.Nil(t, err)
	assert.Equal(t, "SetJohn", v.Name)
	assert.Equal(t, 40, v.Age)
}
