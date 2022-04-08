package mapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	v := struct {
		Name      string `path:"name"`
		Address   string `json:"address,options=[beijing,shanghai]"`
		Age       int    `json:"age"`
		Anonymous bool
	}{
		Name:      "kevin",
		Address:   "shanghai",
		Age:       20,
		Anonymous: true,
	}

	m, err := Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, "kevin", m["path"]["name"])
	assert.Equal(t, "shanghai", m["json"]["address"])
	assert.Equal(t, 20, m["json"]["age"].(int))
	assert.True(t, m[emptyTag]["Anonymous"].(bool))
}

func TestMarshal_BadOptions(t *testing.T) {
	v := struct {
		Name string `json:"name,options"`
	}{
		Name: "kevin",
	}

	_, err := Marshal(v)
	assert.NotNil(t, err)
}

func TestMarshal_NotInOptions(t *testing.T) {
	v := struct {
		Name string `json:"name,options=[a,b]"`
	}{
		Name: "kevin",
	}

	_, err := Marshal(v)
	assert.NotNil(t, err)
}

func TestMarshal_Nested(t *testing.T) {
	type address struct {
		Country string `json:"country"`
		City    string `json:"city"`
	}
	v := struct {
		Name    string  `json:"name,options=[kevin,wan]"`
		Address address `json:"address"`
	}{
		Name: "kevin",
		Address: address{
			Country: "China",
			City:    "Shanghai",
		},
	}

	m, err := Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, "kevin", m["json"]["name"])
	assert.Equal(t, "China", m["json"]["address"].(address).Country)
	assert.Equal(t, "Shanghai", m["json"]["address"].(address).City)
}

func TestMarshal_Slice(t *testing.T) {
	v := struct {
		Name []string `json:"name"`
	}{
		Name: []string{"kevin", "wan"},
	}

	m, err := Marshal(v)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"kevin", "wan"}, m["json"]["name"].([]string))
}

func TestMarshal_SliceNil(t *testing.T) {
	v := struct {
		Name []string `json:"name"`
	}{
		Name: nil,
	}

	_, err := Marshal(v)
	assert.NotNil(t, err)
}
