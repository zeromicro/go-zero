package testjson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	type A struct {
		A  string `json:"a"`
		AA string `json:"aa"`
	}
	type B struct {
		A        // can't be A A, or A `json...`
		B string `json:"b"`
	}
	type C struct {
		A `json:"a"`
		C string `json:"c"`
	}
	a := A{A: "a", AA: "aa"}
	b := B{A: a, B: "b"}
	c := C{A: a, C: "c"}

	bstr, _ := json.Marshal(b)
	cstr, _ := json.Marshal(c)
	assert.Equal(t, `{"a":"a","aa":"aa","b":"b"}`, string(bstr))
	assert.Equal(t, `{"a":{"a":"a","aa":"aa"},"c":"c"}`, string(cstr))
}
