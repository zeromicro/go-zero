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

func TestMarshal_Anonymous(t *testing.T) {
	t.Run("anonymous", func(t *testing.T) {
		type BaseHeader struct {
			Token string `header:"token"`
		}
		v := struct {
			Name    string `json:"name"`
			Address string `json:"address,options=[beijing,shanghai]"`
			Age     int    `json:"age"`
			BaseHeader
		}{
			Name:    "kevin",
			Address: "shanghai",
			Age:     20,
			BaseHeader: BaseHeader{
				Token: "token_xxx",
			},
		}
		m, err := Marshal(v)
		assert.Nil(t, err)
		assert.Equal(t, "kevin", m["json"]["name"])
		assert.Equal(t, "shanghai", m["json"]["address"])
		assert.Equal(t, 20, m["json"]["age"].(int))
		assert.Equal(t, "token_xxx", m["header"]["token"])

		v1 := struct {
			Name    string `json:"name"`
			Address string `json:"address,options=[beijing,shanghai]"`
			Age     int    `json:"age"`
			BaseHeader
		}{
			Name:    "kevin",
			Address: "shanghai",
			Age:     20,
		}
		m1, err1 := Marshal(v1)
		assert.Nil(t, err1)
		assert.Equal(t, "kevin", m1["json"]["name"])
		assert.Equal(t, "shanghai", m1["json"]["address"])
		assert.Equal(t, 20, m1["json"]["age"].(int))

		type AnotherHeader struct {
			Version string `header:"version"`
		}
		v2 := struct {
			Name    string `json:"name"`
			Address string `json:"address,options=[beijing,shanghai]"`
			Age     int    `json:"age"`
			BaseHeader
			AnotherHeader
		}{
			Name:    "kevin",
			Address: "shanghai",
			Age:     20,
			BaseHeader: BaseHeader{
				Token: "token_xxx",
			},
			AnotherHeader: AnotherHeader{
				Version: "v1.0",
			},
		}
		m2, err2 := Marshal(v2)
		assert.Nil(t, err2)
		assert.Equal(t, "kevin", m2["json"]["name"])
		assert.Equal(t, "shanghai", m2["json"]["address"])
		assert.Equal(t, 20, m2["json"]["age"].(int))
		assert.Equal(t, "token_xxx", m2["header"]["token"])
		assert.Equal(t, "v1.0", m2["header"]["version"])

		type PointerHeader struct {
			Ref *string `header:"ref"`
		}
		ref := "reference"
		v3 := struct {
			Name    string `json:"name"`
			Address string `json:"address,options=[beijing,shanghai]"`
			Age     int    `json:"age"`
			PointerHeader
		}{
			Name:    "kevin",
			Address: "shanghai",
			Age:     20,
			PointerHeader: PointerHeader{
				Ref: &ref,
			},
		}
		m3, err3 := Marshal(v3)
		assert.Nil(t, err3)
		assert.Equal(t, "kevin", m3["json"]["name"])
		assert.Equal(t, "shanghai", m3["json"]["address"])
		assert.Equal(t, 20, m3["json"]["age"].(int))
		assert.Equal(t, "reference", *m3["header"]["ref"].(*string))
	})

	t.Run("bad anonymous", func(t *testing.T) {
		type BaseHeader struct {
			Token string `json:"token,options=[a,b]"`
		}

		v := struct {
			Name    string `json:"name"`
			Address string `json:"address,options=[beijing,shanghai]"`
			Age     int    `json:"age"`
			BaseHeader
		}{
			Name:    "kevin",
			Address: "shanghai",
			Age:     20,
			BaseHeader: BaseHeader{
				Token: "c",
			},
		}

		_, err := Marshal(v)
		assert.NotNil(t, err)
	})
}

func TestMarshal_Ptr(t *testing.T) {
	v := &struct {
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

func TestMarshal_OptionalPtr(t *testing.T) {
	var val = 1
	v := struct {
		Age *int `json:"age"`
	}{
		Age: &val,
	}

	m, err := Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, 1, *m["json"]["age"].(*int))
}

func TestMarshal_OptionalPtrNil(t *testing.T) {
	v := struct {
		Age *int `json:"age"`
	}{}

	_, err := Marshal(v)
	assert.NotNil(t, err)
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

func TestMarshal_NotInOptionsOptional(t *testing.T) {
	v := struct {
		Name string `json:"name,options=[a,b],optional"`
	}{}

	_, err := Marshal(v)
	assert.Nil(t, err)
}

func TestMarshal_NotInOptionsOptionalWrongValue(t *testing.T) {
	v := struct {
		Name string `json:"name,options=[a,b],optional"`
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

func TestMarshal_NestedPtr(t *testing.T) {
	type address struct {
		Country string `json:"country"`
		City    string `json:"city"`
	}
	v := struct {
		Name    string   `json:"name,options=[kevin,wan]"`
		Address *address `json:"address"`
	}{
		Name: "kevin",
		Address: &address{
			Country: "China",
			City:    "Shanghai",
		},
	}

	m, err := Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, "kevin", m["json"]["name"])
	assert.Equal(t, "China", m["json"]["address"].(*address).Country)
	assert.Equal(t, "Shanghai", m["json"]["address"].(*address).City)
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

func TestMarshal_Range(t *testing.T) {
	v := struct {
		Int     int     `json:"int,range=[1:3]"`
		Int8    int8    `json:"int8,range=[1:3)"`
		Int16   int16   `json:"int16,range=(1:3]"`
		Int32   int32   `json:"int32,range=(1:3)"`
		Int64   int64   `json:"int64,range=(1:3)"`
		Uint    uint    `json:"uint,range=[1:3]"`
		Uint8   uint8   `json:"uint8,range=[1:3)"`
		Uint16  uint16  `json:"uint16,range=(1:3]"`
		Uint32  uint32  `json:"uint32,range=(1:3)"`
		Uint64  uint64  `json:"uint64,range=(1:3)"`
		Float32 float32 `json:"float32,range=(1:3)"`
		Float64 float64 `json:"float64,range=(1:3)"`
	}{
		Int:     1,
		Int8:    1,
		Int16:   2,
		Int32:   2,
		Int64:   2,
		Uint:    1,
		Uint8:   1,
		Uint16:  2,
		Uint32:  2,
		Uint64:  2,
		Float32: 2,
		Float64: 2,
	}

	m, err := Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, 1, m["json"]["int"].(int))
	assert.Equal(t, int8(1), m["json"]["int8"].(int8))
	assert.Equal(t, int16(2), m["json"]["int16"].(int16))
	assert.Equal(t, int32(2), m["json"]["int32"].(int32))
	assert.Equal(t, int64(2), m["json"]["int64"].(int64))
	assert.Equal(t, uint(1), m["json"]["uint"].(uint))
	assert.Equal(t, uint8(1), m["json"]["uint8"].(uint8))
	assert.Equal(t, uint16(2), m["json"]["uint16"].(uint16))
	assert.Equal(t, uint32(2), m["json"]["uint32"].(uint32))
	assert.Equal(t, uint64(2), m["json"]["uint64"].(uint64))
	assert.Equal(t, float32(2), m["json"]["float32"].(float32))
	assert.Equal(t, float64(2), m["json"]["float64"].(float64))
}

func TestMarshal_RangeOut(t *testing.T) {
	tests := []any{
		struct {
			Int int `json:"int,range=[1:3]"`
		}{
			Int: 4,
		},
		struct {
			Int int `json:"int,range=(1:3]"`
		}{
			Int: 1,
		},
		struct {
			Int int `json:"int,range=[1:3)"`
		}{
			Int: 3,
		},
		struct {
			Int int `json:"int,range=(1:3)"`
		}{
			Int: 3,
		},
		struct {
			Bool bool `json:"bool,range=(1:3)"`
		}{
			Bool: true,
		},
	}

	for _, test := range tests {
		_, err := Marshal(test)
		assert.NotNil(t, err)
	}
}

func TestMarshal_RangeIllegal(t *testing.T) {
	tests := []any{
		struct {
			Int int `json:"int,range=[3:1]"`
		}{
			Int: 2,
		},
		struct {
			Int int `json:"int,range=(3:1]"`
		}{
			Int: 2,
		},
	}

	for _, test := range tests {
		_, err := Marshal(test)
		assert.Equal(t, err, errNumberRange)
	}
}

func TestMarshal_RangeLeftEqualsToRight(t *testing.T) {
	tests := []struct {
		name  string
		value any
		err   error
	}{
		{
			name: "left inclusive, right inclusive",
			value: struct {
				Int int `json:"int,range=[2:2]"`
			}{
				Int: 2,
			},
		},
		{
			name: "left inclusive, right exclusive",
			value: struct {
				Int int `json:"int,range=[2:2)"`
			}{
				Int: 2,
			},
			err: errNumberRange,
		},
		{
			name: "left exclusive, right inclusive",
			value: struct {
				Int int `json:"int,range=(2:2]"`
			}{
				Int: 2,
			},
			err: errNumberRange,
		},
		{
			name: "left exclusive, right exclusive",
			value: struct {
				Int int `json:"int,range=(2:2)"`
			}{
				Int: 2,
			},
			err: errNumberRange,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			_, err := Marshal(test.value)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestMarshal_FromString(t *testing.T) {
	v := struct {
		Age int `json:"age,string"`
	}{
		Age: 10,
	}

	m, err := Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, "10", m["json"]["age"].(string))
}
