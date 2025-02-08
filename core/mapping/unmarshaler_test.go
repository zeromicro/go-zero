package mapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/stringx"
)

// because json.Number doesn't support strconv.ParseUint(...),
// so we only can test to 62 bits.
const maxUintBitsToTest = 62

func TestUnmarshalWithFullNameNotStruct(t *testing.T) {
	var s map[string]any
	content := []byte(`{"name":"xiaoming"}`)
	err := UnmarshalJsonBytes(content, &s)
	assert.Equal(t, errTypeMismatch, err)
}

func TestUnmarshalValueNotSettable(t *testing.T) {
	var s map[string]any
	content := []byte(`{"name":"xiaoming"}`)
	err := UnmarshalJsonBytes(content, s)
	assert.Equal(t, errValueNotSettable, err)
}

func TestUnmarshalWithoutTagName(t *testing.T) {
	type inner struct {
		Optional   bool   `key:",optional"`
		OptionalP  *bool  `key:",optional"`
		OptionalPP **bool `key:",optional"`
	}
	m := map[string]any{
		"Optional":   true,
		"OptionalP":  true,
		"OptionalPP": true,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.True(t, in.Optional)
		assert.True(t, *in.OptionalP)
		assert.True(t, **in.OptionalPP)
	}
}

func TestUnmarshalWithLowerField(t *testing.T) {
	type (
		Lower struct {
			value int `key:"lower"`
		}

		inner struct {
			Lower
			Optional bool `key:",optional"`
		}
	)
	m := map[string]any{
		"Optional": true,
		"lower":    1,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.True(t, in.Optional)
		assert.Equal(t, 0, in.value)
	}
}

func TestUnmarshalWithLowerAnonymousStruct(t *testing.T) {
	type (
		lower struct {
			Value int `key:"lower"`
		}

		inner struct {
			lower
			Optional bool `key:",optional"`
		}
	)
	m := map[string]any{
		"Optional": true,
		"lower":    1,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.True(t, in.Optional)
		assert.Equal(t, 1, in.Value)
	}
}

func TestUnmarshalWithoutTagNameWithCanonicalKey(t *testing.T) {
	type inner struct {
		Name string `key:"name"`
	}
	m := map[string]any{
		"Name": "go-zero",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithCanonicalKeyFunc(func(s string) string {
		first := true
		return strings.Map(func(r rune) rune {
			if first {
				first = false
				return unicode.ToTitle(r)
			}
			return r
		}, s)
	}))
	if assert.NoError(t, unmarshaler.Unmarshal(m, &in)) {
		assert.Equal(t, "go-zero", in.Name)
	}
}

func TestUnmarshalWithoutTagNameWithCanonicalKeyOptionalDep(t *testing.T) {
	type inner struct {
		FirstName string `key:",optional"`
		LastName  string `key:",optional=FirstName"`
	}
	m := map[string]any{
		"firstname": "go",
		"lastname":  "zero",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithCanonicalKeyFunc(func(s string) string {
		return strings.ToLower(s)
	}))
	if assert.NoError(t, unmarshaler.Unmarshal(m, &in)) {
		assert.Equal(t, "go", in.FirstName)
		assert.Equal(t, "zero", in.LastName)
	}
}

func TestUnmarshalBool(t *testing.T) {
	type inner struct {
		True           bool `key:"yes"`
		False          bool `key:"no"`
		TrueFromOne    bool `key:"yesone,string"`
		FalseFromZero  bool `key:"nozero,string"`
		TrueFromTrue   bool `key:"yestrue,string"`
		FalseFromFalse bool `key:"nofalse,string"`
		DefaultTrue    bool `key:"defaulttrue,default=1"`
		Optional       bool `key:"optional,optional"`
	}
	m := map[string]any{
		"yes":     true,
		"no":      false,
		"yesone":  "1",
		"nozero":  "0",
		"yestrue": "true",
		"nofalse": "false",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.True(in.True)
		ast.False(in.False)
		ast.True(in.TrueFromOne)
		ast.False(in.FalseFromZero)
		ast.True(in.TrueFromTrue)
		ast.False(in.FalseFromFalse)
		ast.True(in.DefaultTrue)
	}
}

func TestUnmarshalDuration(t *testing.T) {
	type inner struct {
		Duration       time.Duration   `key:"duration"`
		LessDuration   time.Duration   `key:"less"`
		MoreDuration   time.Duration   `key:"more"`
		PtrDuration    *time.Duration  `key:"ptr"`
		PtrPtrDuration **time.Duration `key:"ptrptr"`
	}
	m := map[string]any{
		"duration": "5s",
		"less":     "100ms",
		"more":     "24h",
		"ptr":      "1h",
		"ptrptr":   "2h",
	}
	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, time.Second*5, in.Duration)
		assert.Equal(t, time.Millisecond*100, in.LessDuration)
		assert.Equal(t, time.Hour*24, in.MoreDuration)
		assert.Equal(t, time.Hour, *in.PtrDuration)
		assert.Equal(t, time.Hour*2, **in.PtrPtrDuration)
	}
}

func TestUnmarshalDurationDefault(t *testing.T) {
	type inner struct {
		Int      int           `key:"int"`
		Duration time.Duration `key:"duration,default=5s"`
	}
	m := map[string]any{
		"int": 5,
	}
	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 5, in.Int)
		assert.Equal(t, time.Second*5, in.Duration)
	}
}

func TestUnmarshalDurationPtr(t *testing.T) {
	type inner struct {
		Duration *time.Duration `key:"duration"`
	}
	m := map[string]any{
		"duration": "5s",
	}
	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, time.Second*5, *in.Duration)
	}
}

func TestUnmarshalDurationPtrDefault(t *testing.T) {
	type inner struct {
		Int      int            `key:"int"`
		Value    *int           `key:",default=5"`
		Duration *time.Duration `key:"duration,default=5s"`
	}
	m := map[string]any{
		"int": 5,
	}
	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 5, in.Int)
		assert.Equal(t, 5, *in.Value)
		assert.Equal(t, time.Second*5, *in.Duration)
	}
}

func TestUnmarshalInt(t *testing.T) {
	type inner struct {
		Int          int   `key:"int"`
		IntFromStr   int   `key:"intstr,string"`
		Int8         int8  `key:"int8"`
		Int8FromStr  int8  `key:"int8str,string"`
		Int16        int16 `key:"int16"`
		Int16FromStr int16 `key:"int16str,string"`
		Int32        int32 `key:"int32"`
		Int32FromStr int32 `key:"int32str,string"`
		Int64        int64 `key:"int64"`
		Int64FromStr int64 `key:"int64str,string"`
		DefaultInt   int64 `key:"defaultint,default=11"`
		Optional     int   `key:"optional,optional"`
		IntOptDef    int   `key:"intopt,optional,default=6"`
	}
	m := map[string]any{
		"int":      1,
		"intstr":   "2",
		"int8":     int8(3),
		"int8str":  "4",
		"int16":    int16(5),
		"int16str": "6",
		"int32":    int32(7),
		"int32str": "8",
		"int64":    int64(9),
		"int64str": "10",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(1, in.Int)
		ast.Equal(2, in.IntFromStr)
		ast.Equal(int8(3), in.Int8)
		ast.Equal(int8(4), in.Int8FromStr)
		ast.Equal(int16(5), in.Int16)
		ast.Equal(int16(6), in.Int16FromStr)
		ast.Equal(int32(7), in.Int32)
		ast.Equal(int32(8), in.Int32FromStr)
		ast.Equal(int64(9), in.Int64)
		ast.Equal(int64(10), in.Int64FromStr)
		ast.Equal(int64(11), in.DefaultInt)
		ast.Equal(6, in.IntOptDef)
	}
}

func TestUnmarshalIntPtr(t *testing.T) {
	type inner struct {
		Int *int `key:"int"`
	}
	m := map[string]any{
		"int": 1,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.NotNil(t, in.Int)
		assert.Equal(t, 1, *in.Int)
	}
}

func TestUnmarshalIntSliceOfPtr(t *testing.T) {
	t.Run("int slice", func(t *testing.T) {
		type inner struct {
			Ints  []*int  `key:"ints"`
			Intps []**int `key:"intps"`
		}
		m := map[string]any{
			"ints":  []int{1, 2, 3},
			"intps": []int{1, 2, 3, 4},
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.NotEmpty(t, in.Ints)
			var ints []int
			for _, i := range in.Ints {
				ints = append(ints, *i)
			}
			assert.EqualValues(t, []int{1, 2, 3}, ints)

			var intps []int
			for _, i := range in.Intps {
				intps = append(intps, **i)
			}
			assert.EqualValues(t, []int{1, 2, 3, 4}, intps)
		}
	})

	t.Run("int slice with error", func(t *testing.T) {
		type inner struct {
			Ints  []*int  `key:"ints"`
			Intps []**int `key:"intps"`
		}
		m := map[string]any{
			"ints":  []any{1, 2, "a"},
			"intps": []int{1, 2, 3, 4},
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int slice with nil element", func(t *testing.T) {
		type inner struct {
			Ints []int `key:"ints"`
		}

		m := map[string]any{
			"ints": []any{nil},
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Empty(t, in.Ints)
		}
	})

	t.Run("int slice with nil", func(t *testing.T) {
		type inner struct {
			Ints []int `key:"ints"`
		}

		m := map[string]any{
			"ints": []any(nil),
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Empty(t, in.Ints)
		}
	})
}

func TestUnmarshalIntWithDefault(t *testing.T) {
	type inner struct {
		Int   int   `key:"int,default=5"`
		Intp  *int  `key:"intp,default=5"`
		Intpp **int `key:"intpp,default=5"`
	}
	m := map[string]any{
		"int":   1,
		"intp":  2,
		"intpp": 3,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 1, in.Int)
		assert.Equal(t, 2, *in.Intp)
		assert.Equal(t, 3, **in.Intpp)
	}
}

func TestUnmarshalIntWithString(t *testing.T) {
	t.Run("int without options", func(t *testing.T) {
		type inner struct {
			Int   int64   `key:"int,string"`
			Intp  *int64  `key:"intp,string"`
			Intpp **int64 `key:"intpp,string"`
		}
		m := map[string]any{
			"int":   json.Number("1"),
			"intp":  json.Number("2"),
			"intpp": json.Number("3"),
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), in.Int)
			assert.Equal(t, int64(2), *in.Intp)
			assert.Equal(t, int64(3), **in.Intpp)
		}
	})

	t.Run("int wrong range", func(t *testing.T) {
		type inner struct {
			Int   int64   `key:"int,string,range=[2:3]"`
			Intp  *int64  `key:"intp,range=[2:3]"`
			Intpp **int64 `key:"intpp,range=[2:3]"`
		}
		m := map[string]any{
			"int":   json.Number("1"),
			"intp":  json.Number("2"),
			"intpp": json.Number("3"),
		}

		var in inner
		assert.ErrorIs(t, UnmarshalKey(m, &in), errNumberRange)
	})

	t.Run("int with wrong type", func(t *testing.T) {
		type (
			myString string

			inner struct {
				Int   int64   `key:"int,string"`
				Intp  *int64  `key:"intp,string"`
				Intpp **int64 `key:"intpp,string"`
			}
		)
		m := map[string]any{
			"int":   myString("1"),
			"intp":  myString("2"),
			"intpp": myString("3"),
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int with ptr", func(t *testing.T) {
		type inner struct {
			Int *int64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("1"),
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), *in.Int)
		}
	})

	t.Run("int with invalid value", func(t *testing.T) {
		type inner struct {
			Int int64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("a"),
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint with invalid value", func(t *testing.T) {
		type inner struct {
			Int uint64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("a"),
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float with invalid value", func(t *testing.T) {
		type inner struct {
			Value float64 `key:"float"`
		}
		m := map[string]any{
			"float": json.Number("a"),
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float with invalid value", func(t *testing.T) {
		type inner struct {
			Value string `key:"value"`
		}
		m := map[string]any{
			"value": json.Number("a"),
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int with ptr of ptr", func(t *testing.T) {
		type inner struct {
			Int **int64 `key:"int"`
		}
		m := map[string]any{
			"int": json.Number("1"),
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), **in.Int)
		}
	})

	t.Run("int with options", func(t *testing.T) {
		type inner struct {
			Int int64 `key:"int,string,options=[0,1]"`
		}
		m := map[string]any{
			"int": json.Number("1"),
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), in.Int)
		}
	})

	t.Run("int with options", func(t *testing.T) {
		type inner struct {
			Int int64 `key:"int,string,options=[0,1]"`
		}
		m := map[string]any{
			"int": nil,
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int with options", func(t *testing.T) {
		type (
			StrType string

			inner struct {
				Int int64 `key:"int,string,options=[0,1]"`
			}
		)
		m := map[string]any{
			"int": StrType("1"),
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("invalid options", func(t *testing.T) {
		type Value struct {
			Name string `key:"name,options="`
		}

		var v Value
		assert.Error(t, UnmarshalKey(emptyMap, &v))
	})
}

func TestUnmarshalInt8WithOverflow(t *testing.T) {
	t.Run("int8 from string", func(t *testing.T) {
		type inner struct {
			Value int8 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int8 from json.Number", func(t *testing.T) {
		type inner struct {
			Value int8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int8 from json.Number", func(t *testing.T) {
		type inner struct {
			Value int8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int8 from int64", func(t *testing.T) {
		type inner struct {
			Value int8 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalInt16WithOverflow(t *testing.T) {
	t.Run("int16 from string", func(t *testing.T) {
		type inner struct {
			Value int16 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int16 from json.Number", func(t *testing.T) {
		type inner struct {
			Value int16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int16 from json.Number", func(t *testing.T) {
		type inner struct {
			Value int16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int16 from int64", func(t *testing.T) {
		type inner struct {
			Value int16 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalInt32WithOverflow(t *testing.T) {
	t.Run("int32 from string", func(t *testing.T) {
		type inner struct {
			Value int32 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int32 from json.Number", func(t *testing.T) {
		type inner struct {
			Value int32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int32 from json.Number", func(t *testing.T) {
		type inner struct {
			Value int32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int32 from int64", func(t *testing.T) {
		type inner struct {
			Value int32 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalInt64WithOverflow(t *testing.T) {
	t.Run("int64 from string", func(t *testing.T) {
		type inner struct {
			Value int64 `key:"int,string"`
		}

		m := map[string]any{
			"int": "18446744073709551616", // overflow, 1 << 64
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("int64 from json.Number", func(t *testing.T) {
		type inner struct {
			Value int64 `key:"int,string"`
		}

		m := map[string]any{
			"int": json.Number("18446744073709551616"), // overflow, 1 << 64
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalUint8WithOverflow(t *testing.T) {
	t.Run("uint8 from string", func(t *testing.T) {
		type inner struct {
			Value uint8 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint8 from json.Number", func(t *testing.T) {
		type inner struct {
			Value uint8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint8 from json.Number with negative", func(t *testing.T) {
		type inner struct {
			Value uint8 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-1"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint8 from int64", func(t *testing.T) {
		type inner struct {
			Value uint8 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalUint16WithOverflow(t *testing.T) {
	t.Run("uint16 from string", func(t *testing.T) {
		type inner struct {
			Value uint16 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint16 from json.Number", func(t *testing.T) {
		type inner struct {
			Value uint16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint16 from json.Number with negative", func(t *testing.T) {
		type inner struct {
			Value uint16 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-1"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint16 from int64", func(t *testing.T) {
		type inner struct {
			Value uint16 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalUint32WithOverflow(t *testing.T) {
	t.Run("uint32 from string", func(t *testing.T) {
		type inner struct {
			Value uint32 `key:"int,string"`
		}

		m := map[string]any{
			"int": "8589934592", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint32 from json.Number", func(t *testing.T) {
		type inner struct {
			Value uint32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("8589934592"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint32 from json.Number with negative", func(t *testing.T) {
		type inner struct {
			Value uint32 `key:"int"`
		}

		m := map[string]any{
			"int": json.Number("-1"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint32 from int64", func(t *testing.T) {
		type inner struct {
			Value uint32 `key:"int"`
		}

		m := map[string]any{
			"int": int64(1) << 36, // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalUint64WithOverflow(t *testing.T) {
	t.Run("uint64 from string", func(t *testing.T) {
		type inner struct {
			Value uint64 `key:"int,string"`
		}

		m := map[string]any{
			"int": "18446744073709551616", // overflow, 1 << 64
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("uint64 from json.Number", func(t *testing.T) {
		type inner struct {
			Value uint64 `key:"int,string"`
		}

		m := map[string]any{
			"int": json.Number("18446744073709551616"), // overflow, 1 << 64
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalFloat32WithOverflow(t *testing.T) {
	t.Run("float32 from string greater than float64", func(t *testing.T) {
		type inner struct {
			Value float32 `key:"float,string"`
		}

		m := map[string]any{
			"float": "1.79769313486231570814527423731704356798070e+309", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from string greater than float32", func(t *testing.T) {
		type inner struct {
			Value float32 `key:"float,string"`
		}

		m := map[string]any{
			"float": "1.79769313486231570814527423731704356798070e+300", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from string less than float32", func(t *testing.T) {
		type inner struct {
			Value float32 `key:"float, string"`
		}

		m := map[string]any{
			"float": "-1.79769313486231570814527423731704356798070e+300", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json.Number greater than float64", func(t *testing.T) {
		type inner struct {
			Value float32 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number("1.79769313486231570814527423731704356798070e+309"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json.Number greater than float32", func(t *testing.T) {
		type inner struct {
			Value float32 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number("1.79769313486231570814527423731704356798070e+300"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json number less than float32", func(t *testing.T) {
		type inner struct {
			Value float32 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number("-1.79769313486231570814527423731704356798070e+300"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalFloat64WithOverflow(t *testing.T) {
	t.Run("float64 from string greater than float64", func(t *testing.T) {
		type inner struct {
			Value float64 `key:"float,string"`
		}

		m := map[string]any{
			"float": "1.79769313486231570814527423731704356798070e+309", // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("float32 from json.Number greater than float64", func(t *testing.T) {
		type inner struct {
			Value float64 `key:"float"`
		}

		m := map[string]any{
			"float": json.Number("1.79769313486231570814527423731704356798070e+309"), // overflow
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalBoolSliceRequired(t *testing.T) {
	type inner struct {
		Bools []bool `key:"bools"`
	}

	var in inner
	assert.NotNil(t, UnmarshalKey(map[string]any{}, &in))
}

func TestUnmarshalBoolSliceNil(t *testing.T) {
	type inner struct {
		Bools []bool `key:"bools,optional"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(map[string]any{}, &in)) {
		assert.Nil(t, in.Bools)
	}
}

func TestUnmarshalBoolSliceNilExplicit(t *testing.T) {
	type inner struct {
		Bools []bool `key:"bools,optional"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"bools": nil,
	}, &in)) {
		assert.Nil(t, in.Bools)
	}
}

func TestUnmarshalBoolSliceEmpty(t *testing.T) {
	type inner struct {
		Bools []bool `key:"bools,optional"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"bools": []bool{},
	}, &in)) {
		assert.Empty(t, in.Bools)
	}
}

func TestUnmarshalBoolSliceWithDefault(t *testing.T) {
	t.Run("slice with default", func(t *testing.T) {
		type inner struct {
			Bools []bool `key:"bools,default=[true,false]"`
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(nil, &in)) {
			assert.ElementsMatch(t, []bool{true, false}, in.Bools)
		}
	})

	t.Run("slice with default error", func(t *testing.T) {
		type inner struct {
			Bools []bool `key:"bools,default=[true,fal]"`
		}

		var in inner
		assert.Error(t, UnmarshalKey(nil, &in))
	})
}

func TestUnmarshalIntSliceWithDefault(t *testing.T) {
	type inner struct {
		Ints []int `key:"ints,default=[1,2,3]"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []int{1, 2, 3}, in.Ints)
	}
}

func TestUnmarshalIntSliceWithDefaultHasSpaces(t *testing.T) {
	type inner struct {
		Ints   []int   `key:"ints,default=[1, 2, 3]"`
		Intps  []*int  `key:"intps,default=[1, 2, 3, 4]"`
		Intpps []**int `key:"intpps,default=[1, 2, 3, 4, 5]"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []int{1, 2, 3}, in.Ints)

		var intps []int
		for _, i := range in.Intps {
			intps = append(intps, *i)
		}
		assert.ElementsMatch(t, []int{1, 2, 3, 4}, intps)

		var intpps []int
		for _, i := range in.Intpps {
			intpps = append(intpps, **i)
		}
		assert.ElementsMatch(t, []int{1, 2, 3, 4, 5}, intpps)
	}
}

func TestUnmarshalFloatSliceWithDefault(t *testing.T) {
	type inner struct {
		Floats []float32 `key:"floats,default=[1.1,2.2,3.3]"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []float32{1.1, 2.2, 3.3}, in.Floats)
	}
}

func TestUnmarshalStringSliceWithDefault(t *testing.T) {
	t.Run("slice with default", func(t *testing.T) {
		type inner struct {
			Strs   []string   `key:"strs,default=[foo,bar,woo]"`
			Strps  []*string  `key:"strs,default=[foo,bar,woo]"`
			Strpps []**string `key:"strs,default=[foo,bar,woo]"`
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(nil, &in)) {
			assert.ElementsMatch(t, []string{"foo", "bar", "woo"}, in.Strs)

			var ss []string
			for _, s := range in.Strps {
				ss = append(ss, *s)
			}
			assert.ElementsMatch(t, []string{"foo", "bar", "woo"}, ss)

			var sss []string
			for _, s := range in.Strpps {
				sss = append(sss, **s)
			}
			assert.ElementsMatch(t, []string{"foo", "bar", "woo"}, sss)
		}
	})

	t.Run("slice with default on errors", func(t *testing.T) {
		type (
			holder struct {
				Chan []chan int
			}

			inner struct {
				Strs []holder `key:"strs,default=[foo,bar,woo]"`
			}
		)

		var in inner
		assert.Error(t, UnmarshalKey(nil, &in))
	})

	t.Run("slice with default on errors", func(t *testing.T) {
		type inner struct {
			Strs []complex64 `key:"strs,default=[foo,bar,woo]"`
		}

		var in inner
		assert.Error(t, UnmarshalKey(nil, &in))
	})
}

func TestUnmarshalStringSliceWithDefaultHasSpaces(t *testing.T) {
	type inner struct {
		Strs []string `key:"strs,default=[foo, bar, woo]"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []string{"foo", "bar", "woo"}, in.Strs)
	}
}

func TestUnmarshalUint(t *testing.T) {
	type inner struct {
		Uint          uint   `key:"uint"`
		UintFromStr   uint   `key:"uintstr,string"`
		Uint8         uint8  `key:"uint8"`
		Uint8FromStr  uint8  `key:"uint8str,string"`
		Uint16        uint16 `key:"uint16"`
		Uint16FromStr uint16 `key:"uint16str,string"`
		Uint32        uint32 `key:"uint32"`
		Uint32FromStr uint32 `key:"uint32str,string"`
		Uint64        uint64 `key:"uint64"`
		Uint64FromStr uint64 `key:"uint64str,string"`
		DefaultUint   uint   `key:"defaultuint,default=11"`
		Optional      uint   `key:"optional,optional"`
	}
	m := map[string]any{
		"uint":      uint(1),
		"uintstr":   "2",
		"uint8":     uint8(3),
		"uint8str":  "4",
		"uint16":    uint16(5),
		"uint16str": "6",
		"uint32":    uint32(7),
		"uint32str": "8",
		"uint64":    uint64(9),
		"uint64str": "10",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(uint(1), in.Uint)
		ast.Equal(uint(2), in.UintFromStr)
		ast.Equal(uint8(3), in.Uint8)
		ast.Equal(uint8(4), in.Uint8FromStr)
		ast.Equal(uint16(5), in.Uint16)
		ast.Equal(uint16(6), in.Uint16FromStr)
		ast.Equal(uint32(7), in.Uint32)
		ast.Equal(uint32(8), in.Uint32FromStr)
		ast.Equal(uint64(9), in.Uint64)
		ast.Equal(uint64(10), in.Uint64FromStr)
		ast.Equal(uint(11), in.DefaultUint)
	}
}

func TestUnmarshalFloat(t *testing.T) {
	type inner struct {
		Float32      float32 `key:"float32"`
		Float32Str   float32 `key:"float32str,string"`
		Float32Num   float32 `key:"float32num"`
		Float64      float64 `key:"float64"`
		Float64Str   float64 `key:"float64str,string"`
		Float64Num   float64 `key:"float64num"`
		DefaultFloat float32 `key:"defaultfloat,default=5.5"`
		Optional     float32 `key:",optional"`
	}
	m := map[string]any{
		"float32":    float32(1.5),
		"float32str": "2.5",
		"float32num": json.Number("2.6"),
		"float64":    3.5,
		"float64str": "4.5",
		"float64num": json.Number("4.6"),
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(float32(1.5), in.Float32)
		ast.Equal(float32(2.5), in.Float32Str)
		ast.Equal(float32(2.6), in.Float32Num)
		ast.Equal(3.5, in.Float64)
		ast.Equal(4.5, in.Float64Str)
		ast.Equal(4.6, in.Float64Num)
		ast.Equal(float32(5.5), in.DefaultFloat)
	}
}

func TestUnmarshalInt64Slice(t *testing.T) {
	var v struct {
		Ages  []int64 `key:"ages"`
		Slice []int64 `key:"slice"`
	}
	m := map[string]any{
		"ages":  []int64{1, 2},
		"slice": []any{},
	}

	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &v)) {
		ast.ElementsMatch([]int64{1, 2}, v.Ages)
		ast.Equal([]int64{}, v.Slice)
	}
}

func TestUnmarshalNullableSlice(t *testing.T) {
	var v struct {
		Ages  []int64 `key:"ages"`
		Slice []int8  `key:"slice"`
	}
	m := map[string]any{
		"ages":  []int64{1, 2},
		"slice": `[null,2]`,
	}

	assert.New(t).Equal(UnmarshalKey(m, &v), errNilSliceElement)
}

func TestUnmarshalWithFloatPtr(t *testing.T) {
	t.Run("*float32", func(t *testing.T) {
		var v struct {
			WeightFloat32 *float32 `key:"weightFloat32,optional"`
		}
		m := map[string]any{
			"weightFloat32": json.Number("3.2"),
		}

		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, float32(3.2), *v.WeightFloat32)
		}
	})

	t.Run("**float32", func(t *testing.T) {
		var v struct {
			WeightFloat32 **float32 `key:"weightFloat32,optional"`
		}
		m := map[string]any{
			"weightFloat32": json.Number("3.2"),
		}

		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, float32(3.2), **v.WeightFloat32)
		}
	})
}

func TestUnmarshalIntSlice(t *testing.T) {
	t.Run("int slice from int", func(t *testing.T) {
		var v struct {
			Ages  []int `key:"ages"`
			Slice []int `key:"slice"`
		}
		m := map[string]any{
			"ages":  []int{1, 2},
			"slice": []any{},
		}

		ast := assert.New(t)
		if ast.NoError(UnmarshalKey(m, &v)) {
			ast.ElementsMatch([]int{1, 2}, v.Ages)
			ast.Equal([]int{}, v.Slice)
		}
	})

	t.Run("int slice from one int", func(t *testing.T) {
		var v struct {
			Ages []int `key:"ages"`
		}
		m := map[string]any{
			"ages": []int{2},
		}

		ast := assert.New(t)
		unmarshaler := NewUnmarshaler(defaultKeyName, WithFromArray())
		if ast.NoError(unmarshaler.Unmarshal(m, &v)) {
			ast.ElementsMatch([]int{2}, v.Ages)
		}
	})

	t.Run("int slice from one int string", func(t *testing.T) {
		var v struct {
			Ages []int `key:"ages"`
		}
		m := map[string]any{
			"ages": []string{"2"},
		}

		ast := assert.New(t)
		unmarshaler := NewUnmarshaler(defaultKeyName, WithFromArray())
		if ast.NoError(unmarshaler.Unmarshal(m, &v)) {
			ast.ElementsMatch([]int{2}, v.Ages)
		}
	})

	t.Run("int slice from one json.Number", func(t *testing.T) {
		var v struct {
			Ages []int `key:"ages"`
		}
		m := map[string]any{
			"ages": []json.Number{"2"},
		}

		ast := assert.New(t)
		unmarshaler := NewUnmarshaler(defaultKeyName, WithFromArray())
		if ast.NoError(unmarshaler.Unmarshal(m, &v)) {
			ast.ElementsMatch([]int{2}, v.Ages)
		}
	})

	t.Run("int slice from one int strings", func(t *testing.T) {
		var v struct {
			Ages []int `key:"ages"`
		}
		m := map[string]any{
			"ages": []string{"1,2"},
		}

		ast := assert.New(t)
		unmarshaler := NewUnmarshaler(defaultKeyName, WithFromArray())
		ast.Error(unmarshaler.Unmarshal(m, &v))
	})
}

func TestUnmarshalString(t *testing.T) {
	type inner struct {
		Name              string `key:"name"`
		NameStr           string `key:"namestr,string"`
		NotPresent        string `key:",optional"`
		NotPresentWithTag string `key:"notpresent,optional"`
		DefaultString     string `key:"defaultstring,default=hello"`
		Optional          string `key:",optional"`
	}
	m := map[string]any{
		"name":    "kevin",
		"namestr": "namewithstring",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal("kevin", in.Name)
		ast.Equal("namewithstring", in.NameStr)
		ast.Empty(in.NotPresent)
		ast.Empty(in.NotPresentWithTag)
		ast.Equal("hello", in.DefaultString)
	}
}

func TestUnmarshalStringWithMissing(t *testing.T) {
	type inner struct {
		Name string `key:"name"`
	}
	m := map[string]any{}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalStringSliceFromString(t *testing.T) {
	t.Run("slice from string", func(t *testing.T) {
		var v struct {
			Names []string `key:"names"`
		}
		m := map[string]any{
			"names": `["first", "second"]`,
		}

		ast := assert.New(t)
		if ast.NoError(UnmarshalKey(m, &v)) {
			ast.Equal(2, len(v.Names))
			ast.Equal("first", v.Names[0])
			ast.Equal("second", v.Names[1])
		}
	})

	t.Run("slice from empty string", func(t *testing.T) {
		var v struct {
			Names []string `key:"names"`
		}
		m := map[string]any{
			"names": []string{""},
		}

		ast := assert.New(t)
		unmarshaler := NewUnmarshaler(defaultKeyName, WithFromArray())
		if ast.NoError(unmarshaler.Unmarshal(m, &v)) {
			ast.ElementsMatch([]string{""}, v.Names)
		}
	})

	t.Run("slice from empty and valid string", func(t *testing.T) {
		var v struct {
			Names []string `key:"names"`
		}
		m := map[string]any{
			"names": []string{","},
		}

		ast := assert.New(t)
		unmarshaler := NewUnmarshaler(defaultKeyName, WithFromArray())
		if ast.NoError(unmarshaler.Unmarshal(m, &v)) {
			ast.ElementsMatch([]string{","}, v.Names)
		}
	})

	t.Run("slice from valid strings with comma", func(t *testing.T) {
		var v struct {
			Names []string `key:"names"`
		}
		m := map[string]any{
			"names": []string{"aa,bb"},
		}

		ast := assert.New(t)
		unmarshaler := NewUnmarshaler(defaultKeyName, WithFromArray())
		if ast.NoError(unmarshaler.Unmarshal(m, &v)) {
			ast.ElementsMatch([]string{"aa,bb"}, v.Names)
		}
	})

	t.Run("slice from string with slice error", func(t *testing.T) {
		var v struct {
			Names []int `key:"names"`
		}
		m := map[string]any{
			"names": `["first", 1]`,
		}

		assert.Error(t, UnmarshalKey(m, &v))
	})

	t.Run("slice from string with error", func(t *testing.T) {
		type myString string

		var v struct {
			Names []string `key:"names"`
		}
		m := map[string]any{
			"names": myString("not a slice"),
		}

		assert.Error(t, UnmarshalKey(m, &v))
	})
}

func TestUnmarshalIntSliceFromString(t *testing.T) {
	var v struct {
		Values []int `key:"values"`
	}
	m := map[string]any{
		"values": `[1, 2]`,
	}

	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &v)) {
		ast.Equal(2, len(v.Values))
		ast.Equal(1, v.Values[0])
		ast.Equal(2, v.Values[1])
	}
}

func TestUnmarshalIntMapFromString(t *testing.T) {
	var v struct {
		Sort map[string]int `key:"sort"`
	}
	m := map[string]any{
		"sort": `{"value":12345,"zeroVal":0,"nullVal":null}`,
	}

	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &v)) {
		ast.Equal(3, len(v.Sort))
		ast.Equal(12345, v.Sort["value"])
		ast.Equal(0, v.Sort["zeroVal"])
		ast.Equal(0, v.Sort["nullVal"])
	}
}

func TestUnmarshalBoolMapFromString(t *testing.T) {
	var v struct {
		Sort map[string]bool `key:"sort"`
	}
	m := map[string]any{
		"sort": `{"value":true,"zeroVal":false,"nullVal":null}`,
	}

	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &v)) {
		ast.Equal(3, len(v.Sort))
		ast.Equal(true, v.Sort["value"])
		ast.Equal(false, v.Sort["zeroVal"])
		ast.Equal(false, v.Sort["nullVal"])
	}
}

type CustomStringer string

type UnsupportedStringer string

func (c CustomStringer) String() string {
	return fmt.Sprintf("{%s}", string(c))
}

func TestUnmarshalStringMapFromStringer(t *testing.T) {
	t.Run("CustomStringer", func(t *testing.T) {
		var v struct {
			Sort map[string]string `key:"sort"`
		}
		m := map[string]any{
			"sort": CustomStringer(`"value":"ascend","emptyStr":""`),
		}

		ast := assert.New(t)
		if ast.NoError(UnmarshalKey(m, &v)) {
			ast.Equal(2, len(v.Sort))
			ast.Equal("ascend", v.Sort["value"])
			ast.Equal("", v.Sort["emptyStr"])
		}
	})

	t.Run("CustomStringer incorrect", func(t *testing.T) {
		var v struct {
			Sort map[string]string `key:"sort"`
		}
		m := map[string]any{
			"sort": CustomStringer(`"value"`),
		}

		assert.Error(t, UnmarshalKey(m, &v))
	})
}

func TestUnmarshalStringMapFromUnsupportedType(t *testing.T) {
	var v struct {
		Sort map[string]string `key:"sort"`
	}
	m := map[string]any{
		"sort": UnsupportedStringer(`{"value":"ascend","emptyStr":""}`),
	}

	ast := assert.New(t)
	ast.Error(UnmarshalKey(m, &v))
}

func TestUnmarshalStringMapFromNotSettableValue(t *testing.T) {
	var v struct {
		sort  map[string]string  `key:"sort"`
		psort *map[string]string `key:"psort"`
	}
	m := map[string]any{
		"sort":  `{"value":"ascend","emptyStr":""}`,
		"psort": `{"value":"ascend","emptyStr":""}`,
	}

	ast := assert.New(t)
	ast.NoError(UnmarshalKey(m, &v))
	assert.Empty(t, v.sort)
	assert.Nil(t, v.psort)
}

func TestUnmarshalStringMapFromString(t *testing.T) {
	var v struct {
		Sort map[string]string `key:"sort"`
	}
	m := map[string]any{
		"sort": `{"value":"ascend","emptyStr":""}`,
	}

	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &v)) {
		ast.Equal(2, len(v.Sort))
		ast.Equal("ascend", v.Sort["value"])
		ast.Equal("", v.Sort["emptyStr"])
	}
}

func TestUnmarshalStructMapFromString(t *testing.T) {
	var v struct {
		Filter map[string]struct {
			Field1 bool     `json:"field1"`
			Field2 int64    `json:"field2,string"`
			Field3 string   `json:"field3"`
			Field4 *string  `json:"field4"`
			Field5 []string `json:"field5"`
		} `key:"filter"`
	}
	m := map[string]any{
		"filter": `{"obj":{"field1":true,"field2":"1573570455447539712","field3":"this is a string",
			"field4":"this is a string pointer","field5":["str1","str2"]}}`,
	}

	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &v)) {
		ast.Equal(1, len(v.Filter))
		ast.NotNil(v.Filter["obj"])
		ast.Equal(true, v.Filter["obj"].Field1)
		ast.Equal(int64(1573570455447539712), v.Filter["obj"].Field2)
		ast.Equal("this is a string", v.Filter["obj"].Field3)
		ast.Equal("this is a string pointer", *v.Filter["obj"].Field4)
		ast.ElementsMatch([]string{"str1", "str2"}, v.Filter["obj"].Field5)
	}
}

func TestUnmarshalStringSliceMapFromString(t *testing.T) {
	var v struct {
		Filter map[string][]string `key:"filter"`
	}
	m := map[string]any{
		"filter": `{"assignType":null,"status":["process","comment"],"rate":[]}`,
	}

	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &v)) {
		ast.Equal(3, len(v.Filter))
		ast.Equal([]string(nil), v.Filter["assignType"])
		ast.Equal(2, len(v.Filter["status"]))
		ast.Equal("process", v.Filter["status"][0])
		ast.Equal("comment", v.Filter["status"][1])
		ast.Equal(0, len(v.Filter["rate"]))
	}
}

func TestUnmarshalStruct(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		type address struct {
			City          string `key:"city"`
			ZipCode       int    `key:"zipcode,string"`
			DefaultString string `key:"defaultstring,default=hello"`
			Optional      string `key:",optional"`
		}
		type inner struct {
			Name      string    `key:"name"`
			Address   address   `key:"address"`
			AddressP  *address  `key:"addressp"`
			AddressPP **address `key:"addresspp"`
		}
		m := map[string]any{
			"name": "kevin",
			"address": map[string]any{
				"city":    "shanghai",
				"zipcode": "200000",
			},
			"addressp": map[string]any{
				"city":    "beijing",
				"zipcode": "300000",
			},
			"addresspp": map[string]any{
				"city":    "guangzhou",
				"zipcode": "400000",
			},
		}

		var in inner
		ast := assert.New(t)
		if ast.NoError(UnmarshalKey(m, &in)) {
			ast.Equal("kevin", in.Name)
			ast.Equal("shanghai", in.Address.City)
			ast.Equal(200000, in.Address.ZipCode)
			ast.Equal("hello", in.AddressP.DefaultString)
			ast.Equal("beijing", in.AddressP.City)
			ast.Equal(300000, in.AddressP.ZipCode)
			ast.Equal("hello", in.AddressP.DefaultString)
			ast.Equal("guangzhou", (*in.AddressPP).City)
			ast.Equal(400000, (*in.AddressPP).ZipCode)
			ast.Equal("hello", (*in.AddressPP).DefaultString)
		}
	})

	t.Run("struct with error", func(t *testing.T) {
		type address struct {
			City          string `key:"city"`
			ZipCode       int    `key:"zipcode,string"`
			DefaultString string `key:"defaultstring,default=hello"`
			Optional      string `key:",optional"`
		}
		type inner struct {
			Name      string    `key:"name"`
			Address   address   `key:"address"`
			AddressP  *address  `key:"addressp"`
			AddressPP **address `key:"addresspp"`
		}
		m := map[string]any{
			"name": "kevin",
			"address": map[string]any{
				"city":    "shanghai",
				"zipcode": "200000",
			},
			"addressp": map[string]any{
				"city":    "beijing",
				"zipcode": "300000",
			},
			"addresspp": map[string]any{
				"city":    "guangzhou",
				"zipcode": "a",
			},
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalStructOptionalDepends(t *testing.T) {
	type address struct {
		City            string `key:"city"`
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=Optional"`
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	tests := []struct {
		input map[string]string
		pass  bool
	}{
		{
			pass: true,
		},
		{
			input: map[string]string{
				"OptionalDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optional": "a",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optional":        "a",
				"OptionalDepends": "b",
			},
			pass: true,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			m := map[string]any{
				"name": "kevin",
				"address": map[string]any{
					"city": "shanghai",
				},
			}
			for k, v := range test.input {
				m["address"].(map[string]any)[k] = v
			}

			var in inner
			ast := assert.New(t)
			if test.pass {
				if ast.NoError(UnmarshalKey(m, &in)) {
					ast.Equal("kevin", in.Name)
					ast.Equal("shanghai", in.Address.City)
					ast.Equal(test.input["Optional"], in.Address.Optional)
					ast.Equal(test.input["OptionalDepends"], in.Address.OptionalDepends)
				}
			} else {
				ast.Error(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestUnmarshalStructOptionalDependsNot(t *testing.T) {
	type address struct {
		City            string `key:"city"`
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=!Optional"`
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	tests := []struct {
		input map[string]string
		pass  bool
	}{
		{
			input: map[string]string{},
			pass:  false,
		},
		{
			input: map[string]string{
				"Optional":        "a",
				"OptionalDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optional": "a",
			},
			pass: true,
		},
		{
			input: map[string]string{
				"OptionalDepends": "b",
			},
			pass: true,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			m := map[string]any{
				"name": "kevin",
				"address": map[string]any{
					"city": "shanghai",
				},
			}
			for k, v := range test.input {
				m["address"].(map[string]any)[k] = v
			}

			var in inner
			ast := assert.New(t)
			if test.pass {
				if ast.NoError(UnmarshalKey(m, &in)) {
					ast.Equal("kevin", in.Name)
					ast.Equal("shanghai", in.Address.City)
					ast.Equal(test.input["Optional"], in.Address.Optional)
					ast.Equal(test.input["OptionalDepends"], in.Address.OptionalDepends)
				}
			} else {
				ast.Error(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestUnmarshalStructOptionalDependsNotErrorDetails(t *testing.T) {
	t.Run("mutal optionals", func(t *testing.T) {
		type address struct {
			Optional        string `key:",optional"`
			OptionalDepends string `key:",optional=!Optional"`
		}
		type inner struct {
			Name    string  `key:"name"`
			Address address `key:"address"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("with default", func(t *testing.T) {
		type address struct {
			Optional        string `key:",optional"`
			OptionalDepends string `key:",default=value,optional"`
		}
		type inner struct {
			Name    string  `key:"name"`
			Address address `key:"address"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, "kevin", in.Name)
			assert.Equal(t, "value", in.Address.OptionalDepends)
		}
	})
}

func TestUnmarshalStructOptionalDependsNotNested(t *testing.T) {
	t.Run("mutal optionals", func(t *testing.T) {
		type address struct {
			Optional        string `key:",optional"`
			OptionalDepends string `key:",optional=!Optional"`
		}
		type combo struct {
			Name    string  `key:"name,optional"`
			Address address `key:"address"`
		}
		type inner struct {
			Name  string `key:"name"`
			Combo combo  `key:"combo"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("bad format", func(t *testing.T) {
		type address struct {
			Optional        string `key:",optional"`
			OptionalDepends string `key:",optional=!Optional=abcd"`
		}
		type combo struct {
			Name    string  `key:"name,optional"`
			Address address `key:"address"`
		}
		type inner struct {
			Name  string `key:"name"`
			Combo combo  `key:"combo"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})

	t.Run("invalid option", func(t *testing.T) {
		type address struct {
			Optional        string `key:",optional"`
			OptionalDepends string `key:",opt=abcd"`
		}
		type combo struct {
			Name    string  `key:"name,optional"`
			Address address `key:"address"`
		}
		type inner struct {
			Name  string `key:"name"`
			Combo combo  `key:"combo"`
		}

		m := map[string]any{
			"name": "kevin",
		}

		var in inner
		assert.Error(t, UnmarshalKey(m, &in))
	})
}

func TestUnmarshalStructOptionalNestedDifferentKey(t *testing.T) {
	type address struct {
		Optional        string `dkey:",optional"`
		OptionalDepends string `key:",optional"`
	}
	type combo struct {
		Name    string  `key:"name,optional"`
		Address address `key:"address"`
	}
	type inner struct {
		Name  string `key:"name"`
		Combo combo  `key:"combo"`
	}

	m := map[string]any{
		"name": "kevin",
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalStructOptionalDependsNotEnoughValue(t *testing.T) {
	type address struct {
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=!"`
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name":    "kevin",
		"address": map[string]any{},
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalStructOptionalDependsMoreValues(t *testing.T) {
	type address struct {
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=a=b"`
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name":    "kevin",
		"address": map[string]any{},
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalStructMissing(t *testing.T) {
	type address struct {
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=a=b"`
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name": "kevin",
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalNestedStructMissing(t *testing.T) {
	type mostInner struct {
		Name string `key:"name"`
	}
	type address struct {
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=a=b"`
		MostInner       mostInner
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]any{
		"name":    "kevin",
		"address": map[string]any{},
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalAnonymousStructOptionalDepends(t *testing.T) {
	type AnonAddress struct {
		City            string `key:"city"`
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=Optional"`
	}
	type inner struct {
		Name string `key:"name"`
		AnonAddress
	}

	tests := []struct {
		input map[string]string
		pass  bool
	}{
		{
			pass: true,
		},
		{
			input: map[string]string{
				"OptionalDepends": "b",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optional": "a",
			},
			pass: false,
		},
		{
			input: map[string]string{
				"Optional":        "a",
				"OptionalDepends": "b",
			},
			pass: true,
		},
	}

	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			m := map[string]any{
				"name": "kevin",
				"city": "shanghai",
			}
			for k, v := range test.input {
				m[k] = v
			}

			var in inner
			ast := assert.New(t)
			if test.pass {
				if ast.NoError(UnmarshalKey(m, &in)) {
					ast.Equal("kevin", in.Name)
					ast.Equal("shanghai", in.City)
					ast.Equal(test.input["Optional"], in.Optional)
					ast.Equal(test.input["OptionalDepends"], in.OptionalDepends)
				}
			} else {
				ast.Error(UnmarshalKey(m, &in))
			}
		})
	}
}

func TestUnmarshalStructPtr(t *testing.T) {
	type address struct {
		City          string `key:"city"`
		ZipCode       int    `key:"zipcode,string"`
		DefaultString string `key:"defaultstring,default=hello"`
		Optional      string `key:",optional"`
	}
	type inner struct {
		Name    string   `key:"name"`
		Address *address `key:"address"`
	}
	m := map[string]any{
		"name": "kevin",
		"address": map[string]any{
			"city":    "shanghai",
			"zipcode": "200000",
		},
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal("kevin", in.Name)
		ast.Equal("shanghai", in.Address.City)
		ast.Equal(200000, in.Address.ZipCode)
		ast.Equal("hello", in.Address.DefaultString)
	}
}

func TestUnmarshalWithStringIgnored(t *testing.T) {
	type inner struct {
		True    bool    `key:"yes"`
		False   bool    `key:"no"`
		Int     int     `key:"int"`
		Int8    int8    `key:"int8"`
		Int16   int16   `key:"int16"`
		Int32   int32   `key:"int32"`
		Int64   int64   `key:"int64"`
		Uint    uint    `key:"uint"`
		Uint8   uint8   `key:"uint8"`
		Uint16  uint16  `key:"uint16"`
		Uint32  uint32  `key:"uint32"`
		Uint64  uint64  `key:"uint64"`
		Float32 float32 `key:"float32"`
		Float64 float64 `key:"float64"`
	}
	m := map[string]any{
		"yes":     "1",
		"no":      "0",
		"int":     "1",
		"int8":    "3",
		"int16":   "5",
		"int32":   "7",
		"int64":   "9",
		"uint":    "1",
		"uint8":   "3",
		"uint16":  "5",
		"uint32":  "7",
		"uint64":  "9",
		"float32": "1.5",
		"float64": "3.5",
	}

	var in inner
	um := NewUnmarshaler("key", WithStringValues())
	ast := assert.New(t)
	if ast.NoError(um.Unmarshal(m, &in)) {
		ast.True(in.True)
		ast.False(in.False)
		ast.Equal(1, in.Int)
		ast.Equal(int8(3), in.Int8)
		ast.Equal(int16(5), in.Int16)
		ast.Equal(int32(7), in.Int32)
		ast.Equal(int64(9), in.Int64)
		ast.Equal(uint(1), in.Uint)
		ast.Equal(uint8(3), in.Uint8)
		ast.Equal(uint16(5), in.Uint16)
		ast.Equal(uint32(7), in.Uint32)
		ast.Equal(uint64(9), in.Uint64)
		ast.Equal(float32(1.5), in.Float32)
		ast.Equal(3.5, in.Float64)
	}
}

func TestUnmarshalJsonNumberInt64(t *testing.T) {
	for i := 0; i <= maxUintBitsToTest; i++ {
		var intValue int64 = 1 << uint(i)
		strValue := strconv.FormatInt(intValue, 10)
		number := json.Number(strValue)
		m := map[string]any{
			"ID": number,
		}
		var v struct {
			ID int64
		}
		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, intValue, v.ID)
		}
	}
}

func TestUnmarshalJsonNumberUint64(t *testing.T) {
	for i := 0; i <= maxUintBitsToTest; i++ {
		var intValue uint64 = 1 << uint(i)
		strValue := strconv.FormatUint(intValue, 10)
		number := json.Number(strValue)
		m := map[string]any{
			"ID": number,
		}
		var v struct {
			ID uint64
		}
		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, intValue, v.ID)
		}
	}
}

func TestUnmarshalJsonNumberUint64Ptr(t *testing.T) {
	for i := 0; i <= maxUintBitsToTest; i++ {
		var intValue uint64 = 1 << uint(i)
		strValue := strconv.FormatUint(intValue, 10)
		number := json.Number(strValue)
		m := map[string]any{
			"ID": number,
		}
		var v struct {
			ID *uint64
		}
		ast := assert.New(t)
		if ast.NoError(UnmarshalKey(m, &v)) {
			ast.NotNil(v.ID)
			ast.Equal(intValue, *v.ID)
		}
	}
}

func TestUnmarshalMapOfInt(t *testing.T) {
	m := map[string]any{
		"Ids": map[string]bool{"first": true},
	}
	var v struct {
		Ids map[string]bool
	}
	if assert.NoError(t, UnmarshalKey(m, &v)) {
		assert.True(t, v.Ids["first"])
	}
}

func TestUnmarshalMapOfStruct(t *testing.T) {
	t.Run("map of struct with error", func(t *testing.T) {
		m := map[string]any{
			"Ids": map[string]any{"first": "second"},
		}
		var v struct {
			Ids map[string]struct {
				Name string
			}
		}
		assert.Error(t, UnmarshalKey(m, &v))
	})

	t.Run("map of struct", func(t *testing.T) {
		m := map[string]any{
			"Ids": map[string]any{
				"foo": map[string]any{"Name": "foo"},
			},
		}
		var v struct {
			Ids map[string]struct {
				Name string
			}
		}
		if assert.NoError(t, UnmarshalKey(m, &v)) {
			assert.Equal(t, "foo", v.Ids["foo"].Name)
		}
	})

	t.Run("map of struct error", func(t *testing.T) {
		m := map[string]any{
			"Ids": map[string]any{
				"foo": map[string]any{"name": "foo"},
			},
		}
		var v struct {
			Ids map[string]struct {
				Name string
			}
		}
		assert.Error(t, UnmarshalKey(m, &v))
	})
}

func TestUnmarshalSlice(t *testing.T) {
	t.Run("slice of string", func(t *testing.T) {
		m := map[string]any{
			"Ids": []any{"first", "second"},
		}
		var v struct {
			Ids []string
		}
		ast := assert.New(t)
		if ast.NoError(UnmarshalKey(m, &v)) {
			ast.Equal(2, len(v.Ids))
			ast.Equal("first", v.Ids[0])
			ast.Equal("second", v.Ids[1])
		}
	})

	t.Run("slice with type mismatch", func(t *testing.T) {
		var v struct {
			Ids string
		}
		assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal([]any{1, 2}, &v))
	})

	t.Run("slice", func(t *testing.T) {
		var v []int
		ast := assert.New(t)
		if ast.NoError(NewUnmarshaler(jsonTagKey).Unmarshal([]any{1, 2}, &v)) {
			ast.Equal(2, len(v))
			ast.Equal(1, v[0])
			ast.Equal(2, v[1])
		}
	})

	t.Run("slice with unsupported type", func(t *testing.T) {
		var v int
		assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(1, &v))
	})
}

func TestUnmarshalSliceOfStruct(t *testing.T) {
	t.Run("slice of struct", func(t *testing.T) {
		m := map[string]any{
			"Ids": []map[string]any{
				{
					"First":  1,
					"Second": 2,
				},
			},
		}
		var v struct {
			Ids []struct {
				First  int
				Second int
			}
		}
		ast := assert.New(t)
		if ast.NoError(UnmarshalKey(m, &v)) {
			ast.Equal(1, len(v.Ids))
			ast.Equal(1, v.Ids[0].First)
			ast.Equal(2, v.Ids[0].Second)
		}
	})

	t.Run("slice of struct", func(t *testing.T) {
		m := map[string]any{
			"Ids": []map[string]any{
				{
					"First":  "a",
					"Second": 2,
				},
			},
		}
		var v struct {
			Ids []struct {
				First  int
				Second int
			}
		}
		assert.Error(t, UnmarshalKey(m, &v))
	})
}

func TestUnmarshalWithStringOptionsCorrect(t *testing.T) {
	type inner struct {
		Value   string `key:"value,options=first|second"`
		Foo     string `key:"foo,options=[bar,baz]"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"foo":     "bar",
		"correct": "2",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal("first", in.Value)
		ast.Equal("bar", in.Foo)
		ast.Equal("2", in.Correct)
	}
}

func TestUnmarshalOptionsOptional(t *testing.T) {
	type inner struct {
		Value         string `key:"value,options=first|second,optional"`
		OptionalValue string `key:"optional_value,options=first|second,optional"`
		Foo           string `key:"foo,options=[bar,baz]"`
		Correct       string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"foo":     "bar",
		"correct": "2",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal("first", in.Value)
		ast.Equal("", in.OptionalValue)
		ast.Equal("bar", in.Foo)
		ast.Equal("2", in.Correct)
	}
}

func TestUnmarshalOptionsOptionalWrongValue(t *testing.T) {
	type inner struct {
		Value         string `key:"value,options=first|second,optional"`
		OptionalValue string `key:"optional_value,options=first|second,optional"`
		WrongValue    string `key:"wrong_value,options=first|second,optional"`
	}
	m := map[string]any{
		"value":       "first",
		"wrong_value": "third",
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalOptionsMissingValues(t *testing.T) {
	type inner struct {
		Value string `key:"value,options"`
	}
	m := map[string]any{
		"value": "first",
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalStringOptionsWithStringOptionsNotString(t *testing.T) {
	type inner struct {
		Value   string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"correct": 2,
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	assert.Error(t, unmarshaler.Unmarshal(m, &in))
}

func TestUnmarshalStringOptionsWithStringOptions(t *testing.T) {
	type inner struct {
		Value   string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"correct": "2",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	ast := assert.New(t)
	if ast.NoError(unmarshaler.Unmarshal(m, &in)) {
		ast.Equal("first", in.Value)
		ast.Equal("2", in.Correct)
	}
}

func TestUnmarshalStringOptionsWithStringOptionsPtr(t *testing.T) {
	type inner struct {
		Value   *string  `key:"value,options=first|second"`
		ValueP  **string `key:"valuep,options=first|second"`
		Correct *int     `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "first",
		"valuep":  "second",
		"correct": "2",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	ast := assert.New(t)
	if ast.NoError(unmarshaler.Unmarshal(m, &in)) {
		ast.True(*in.Value == "first")
		ast.True(**in.ValueP == "second")
		ast.True(*in.Correct == 2)
	}
}

func TestUnmarshalStringOptionsWithStringOptionsIncorrect(t *testing.T) {
	type inner struct {
		Value   string `key:"value,options=first|second"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "third",
		"correct": "2",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	assert.Error(t, unmarshaler.Unmarshal(m, &in))
}

func TestUnmarshalStringOptionsWithStringOptionsIncorrectGrouped(t *testing.T) {
	type inner struct {
		Value   string `key:"value,options=[first,second]"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]any{
		"value":   "third",
		"correct": "2",
	}

	var in inner
	unmarshaler := NewUnmarshaler(defaultKeyName, WithStringValues())
	assert.Error(t, unmarshaler.Unmarshal(m, &in))
}

func TestUnmarshalWithStringOptionsIncorrect(t *testing.T) {
	type inner struct {
		Value     string `key:"value,options=first|second"`
		Incorrect string `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": "3",
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalWithIntOptionsCorrect(t *testing.T) {
	type inner struct {
		Value  string `key:"value,options=first|second"`
		Number int    `key:"number,options=1|2"`
	}
	m := map[string]any{
		"value":  "first",
		"number": 2,
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal("first", in.Value)
		ast.Equal(2, in.Number)
	}
}

func TestUnmarshalWithIntOptionsCorrectPtr(t *testing.T) {
	type inner struct {
		Value  *string `key:"value,options=first|second"`
		Number *int    `key:"number,options=1|2"`
	}
	m := map[string]any{
		"value":  "first",
		"number": 2,
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.True(*in.Value == "first")
		ast.True(*in.Number == 2)
	}
}

func TestUnmarshalWithIntOptionsIncorrect(t *testing.T) {
	type inner struct {
		Value     string `key:"value,options=first|second"`
		Incorrect int    `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": 3,
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalWithJsonNumberOptionsIncorrect(t *testing.T) {
	type inner struct {
		Value     string `key:"value,options=first|second"`
		Incorrect int    `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": json.Number("3"),
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshaler_UnmarshalIntOptions(t *testing.T) {
	var val struct {
		Sex int `json:"sex,options=0|1"`
	}
	input := []byte(`{"sex": 2}`)
	assert.Error(t, UnmarshalJsonBytes(input, &val))
}

func TestUnmarshalWithUintOptionsCorrect(t *testing.T) {
	type inner struct {
		Value  string `key:"value,options=first|second"`
		Number uint   `key:"number,options=1|2"`
	}
	m := map[string]any{
		"value":  "first",
		"number": uint(2),
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal("first", in.Value)
		ast.Equal(uint(2), in.Number)
	}
}

func TestUnmarshalWithUintOptionsIncorrect(t *testing.T) {
	type inner struct {
		Value     string `key:"value,options=first|second"`
		Incorrect uint   `key:"incorrect,options=1|2"`
	}
	m := map[string]any{
		"value":     "first",
		"incorrect": uint(3),
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalWithOptionsAndDefault(t *testing.T) {
	type inner struct {
		Value string `key:"value,options=first|second|third,default=second"`
	}
	m := map[string]any{}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, "second", in.Value)
	}
}

func TestUnmarshalWithOptionsAndSet(t *testing.T) {
	type inner struct {
		Value string `key:"value,options=first|second|third,default=second"`
	}
	m := map[string]any{
		"value": "first",
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, "first", in.Value)
	}
}

func TestUnmarshalNestedKey(t *testing.T) {
	var c struct {
		ID int `json:"Persons.first.ID"`
	}
	m := map[string]any{
		"Persons": map[string]any{
			"first": map[string]any{
				"ID": 1,
			},
		},
	}

	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
		assert.Equal(t, 1, c.ID)
	}
}

func TestUnmarhsalNestedKeyArray(t *testing.T) {
	var c struct {
		First []struct {
			ID int
		} `json:"Persons.first"`
	}
	m := map[string]any{
		"Persons": map[string]any{
			"first": []map[string]any{
				{"ID": 1},
				{"ID": 2},
			},
		},
	}

	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
		assert.Equal(t, 2, len(c.First))
		assert.Equal(t, 1, c.First[0].ID)
	}
}

func TestUnmarshalAnonymousOptionalRequiredProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousOptionalRequiredMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Value) == 0)
	}
}

func TestUnmarshalAnonymousOptionalOptionalProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousOptionalOptionalMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Value) == 0)
	}
}

func TestUnmarshalAnonymousOptionalRequiredBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousOptionalRequiredOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestUnmarshalAnonymousOptionalRequiredBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.True(t, len(b.Value) == 0)
	}
}

func TestUnmarshalAnonymousOptionalOneRequiredOneOptionalBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousOptionalOneRequiredOneOptionalBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.True(t, len(b.Value) == 0)
	}
}

func TestUnmarshalAnonymousOptionalOneRequiredOneOptionalRequiredProvidedOptionalMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousOptionalOneRequiredOneOptionalRequiredMissedOptionalProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"n": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestUnmarshalAnonymousOptionalBothOptionalBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousOptionalBothOptionalOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousOptionalBothOptionalBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo `json:",optional"`
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.True(t, len(b.Value) == 0)
	}
}

func TestUnmarshalAnonymousRequiredProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousRequiredMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestUnmarshalAnonymousOptionalProvided(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousOptionalMissed(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Value) == 0)
	}
}

func TestUnmarshalAnonymousRequiredBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousRequiredOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestUnmarshalAnonymousRequiredBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestUnmarshalAnonymousOneRequiredOneOptionalBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousOneRequiredOneOptionalBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestUnmarshalAnonymousOneRequiredOneOptionalRequiredProvidedOptionalMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousOneRequiredOneOptionalRequiredMissedOptionalProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"n": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestUnmarshalAnonymousBothOptionalBothProvided(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "kevin", b.Name)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousBothOptionalOneProvidedOneMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.Equal(t, "anything", b.Value)
	}
}

func TestUnmarshalAnonymousBothOptionalBothMissed(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n,optional"`
			Value string `json:"v,optional"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.True(t, len(b.Name) == 0)
		assert.True(t, len(b.Value) == 0)
	}
}

func TestUnmarshalAnonymousWrappedToMuch(t *testing.T) {
	type (
		Foo struct {
			Name  string `json:"n"`
			Value string `json:"v"`
		}

		Bar struct {
			Foo
		}
	)
	m := map[string]any{
		"Foo": map[string]any{
			"n": "name",
			"v": "anything",
		},
	}

	var b Bar
	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b))
}

func TestUnmarshalWrappedObject(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Inner Foo
		}
	)
	m := map[string]any{
		"Inner": map[string]any{
			"v": "anything",
		},
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Inner.Value)
	}
}

func TestUnmarshalWrappedObjectOptional(t *testing.T) {
	type (
		Foo struct {
			Hosts []string
			Key   string
		}

		Bar struct {
			Inner Foo `json:",optional"`
			Name  string
		}
	)
	m := map[string]any{
		"Name": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Name)
	}
}

func TestUnmarshalWrappedObjectOptionalFilled(t *testing.T) {
	type (
		Foo struct {
			Hosts []string
			Key   string
		}

		Bar struct {
			Inner Foo `json:",optional"`
			Name  string
		}
	)
	hosts := []string{"1", "2"}
	m := map[string]any{
		"Inner": map[string]any{
			"Hosts": hosts,
			"Key":   "key",
		},
		"Name": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.EqualValues(t, hosts, b.Inner.Hosts)
		assert.Equal(t, "key", b.Inner.Key)
		assert.Equal(t, "anything", b.Name)
	}
}

func TestUnmarshalWrappedNamedObjectOptional(t *testing.T) {
	type (
		Foo struct {
			Host string
			Key  string
		}

		Bar struct {
			Inner Foo `json:",optional"`
			Name  string
		}
	)
	m := map[string]any{
		"Inner": map[string]any{
			"Host": "thehost",
			"Key":  "thekey",
		},
		"Name": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "thehost", b.Inner.Host)
		assert.Equal(t, "thekey", b.Inner.Key)
		assert.Equal(t, "anything", b.Name)
	}
}

func TestUnmarshalWrappedObjectNamedPtr(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Inner *Foo `json:"foo,optional"`
		}
	)
	m := map[string]any{
		"foo": map[string]any{
			"v": "anything",
		},
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Inner.Value)
	}
}

func TestUnmarshalWrappedObjectPtr(t *testing.T) {
	type (
		Foo struct {
			Value string `json:"v"`
		}

		Bar struct {
			Inner *Foo
		}
	)
	m := map[string]any{
		"Inner": map[string]any{
			"v": "anything",
		},
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Inner.Value)
	}
}

func TestUnmarshalInt2String(t *testing.T) {
	type inner struct {
		Int string `key:"int"`
	}
	m := map[string]any{
		"int": 123,
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalZeroValues(t *testing.T) {
	type inner struct {
		False  bool   `key:"no"`
		Int    int    `key:"int"`
		String string `key:"string"`
	}
	m := map[string]any{
		"no":     false,
		"int":    0,
		"string": "",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.False(in.False)
		ast.Equal(0, in.Int)
		ast.Equal("", in.String)
	}
}

func TestUnmarshalUsingDifferentKeys(t *testing.T) {
	type inner struct {
		False  bool   `key:"no"`
		Int    int    `key:"int"`
		String string `bson:"string"`
	}
	m := map[string]any{
		"no":     false,
		"int":    9,
		"string": "value",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.False(in.False)
		ast.Equal(9, in.Int)
		ast.True(len(in.String) == 0)
	}
}

func TestUnmarshalNumberRangeInt(t *testing.T) {
	type inner struct {
		Value1  int    `key:"value1,range=[1:]"`
		Value2  int8   `key:"value2,range=[1:5]"`
		Value3  int16  `key:"value3,range=[1:5]"`
		Value4  int32  `key:"value4,range=[1:5]"`
		Value5  int64  `key:"value5,range=[1:5]"`
		Value6  uint   `key:"value6,range=[:5]"`
		Value8  uint8  `key:"value8,range=[1:5],string"`
		Value9  uint16 `key:"value9,range=[1:5],string"`
		Value10 uint32 `key:"value10,range=[1:5],string"`
		Value11 uint64 `key:"value11,range=[1:5],string"`
	}
	m := map[string]any{
		"value1":  10,
		"value2":  int8(1),
		"value3":  int16(2),
		"value4":  int32(4),
		"value5":  int64(5),
		"value6":  uint(0),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(10, in.Value1)
		ast.Equal(int8(1), in.Value2)
		ast.Equal(int16(2), in.Value3)
		ast.Equal(int32(4), in.Value4)
		ast.Equal(int64(5), in.Value5)
		ast.Equal(uint(0), in.Value6)
		ast.Equal(uint8(1), in.Value8)
		ast.Equal(uint16(2), in.Value9)
		ast.Equal(uint32(4), in.Value10)
		ast.Equal(uint64(5), in.Value11)
	}
}

func TestUnmarshalNumberRangeJsonNumber(t *testing.T) {
	type inner struct {
		Value3 uint   `key:"value3,range=(1:5]"`
		Value4 uint8  `key:"value4,range=(1:5]"`
		Value5 uint16 `key:"value5,range=(1:5]"`
	}
	m := map[string]any{
		"value3": json.Number("2"),
		"value4": json.Number("4"),
		"value5": json.Number("5"),
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(uint(2), in.Value3)
		ast.Equal(uint8(4), in.Value4)
		ast.Equal(uint16(5), in.Value5)
	}

	type inner1 struct {
		Value int `key:"value,range=(1:5]"`
	}
	m = map[string]any{
		"value": json.Number("a"),
	}

	var in1 inner1
	ast.Error(UnmarshalKey(m, &in1))
}

func TestUnmarshalNumberRangeIntLeftExclude(t *testing.T) {
	type inner struct {
		Value3  uint   `key:"value3,range=(1:5]"`
		Value4  uint32 `key:"value4,default=4,range=(1:5]"`
		Value5  uint64 `key:"value5,range=(1:5]"`
		Value9  int    `key:"value9,range=(1:5],string"`
		Value10 int    `key:"value10,range=(1:5],string"`
		Value11 int    `key:"value11,range=(1:5],string"`
	}
	m := map[string]any{
		"value3":  uint(2),
		"value4":  uint32(4),
		"value5":  uint64(5),
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(uint(2), in.Value3)
		ast.Equal(uint32(4), in.Value4)
		ast.Equal(uint64(5), in.Value5)
		ast.Equal(2, in.Value9)
		ast.Equal(4, in.Value10)
		ast.Equal(5, in.Value11)
	}
}

func TestUnmarshalNumberRangeIntRightExclude(t *testing.T) {
	type inner struct {
		Value2  uint   `key:"value2,range=[1:5)"`
		Value3  uint8  `key:"value3,range=[1:5)"`
		Value4  uint16 `key:"value4,range=[1:5)"`
		Value8  int    `key:"value8,range=[1:5),string"`
		Value9  int    `key:"value9,range=[1:5),string"`
		Value10 int    `key:"value10,range=[1:5),string"`
	}
	m := map[string]any{
		"value2":  uint(1),
		"value3":  uint8(2),
		"value4":  uint16(4),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(uint(1), in.Value2)
		ast.Equal(uint8(2), in.Value3)
		ast.Equal(uint16(4), in.Value4)
		ast.Equal(1, in.Value8)
		ast.Equal(2, in.Value9)
		ast.Equal(4, in.Value10)
	}
}

func TestUnmarshalNumberRangeIntExclude(t *testing.T) {
	type inner struct {
		Value3  int `key:"value3,range=(1:5)"`
		Value4  int `key:"value4,range=(1:5)"`
		Value9  int `key:"value9,range=(1:5),string"`
		Value10 int `key:"value10,range=(1:5),string"`
	}
	m := map[string]any{
		"value3":  2,
		"value4":  4,
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(2, in.Value3)
		ast.Equal(4, in.Value4)
		ast.Equal(2, in.Value9)
		ast.Equal(4, in.Value10)
	}
}

func TestUnmarshalNumberRangeIntOutOfRange(t *testing.T) {
	type inner1 struct {
		Value int64 `key:"value,default=3,range=(1:5)"`
	}

	var in1 inner1
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(1),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(0),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(5),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": json.Number("6"),
	}, &in1))

	type inner2 struct {
		Value int64 `key:"value,optional,range=[1:5)"`
	}

	var in2 inner2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(0),
	}, &in2))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(5),
	}, &in2))

	type inner3 struct {
		Value int64 `key:"value,range=(1:5]"`
	}

	var in3 inner3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(1),
	}, &in3))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(6),
	}, &in3))

	type inner4 struct {
		Value int64 `key:"value,range=[1:5]"`
	}

	var in4 inner4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(0),
	}, &in4))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": int64(6),
	}, &in4))
}

func TestUnmarshalNumberRangeFloat(t *testing.T) {
	type inner struct {
		Value2  float32 `key:"value2,range=[1:5]"`
		Value3  float32 `key:"value3,range=[1:5]"`
		Value4  float64 `key:"value4,range=[1:5]"`
		Value5  float64 `key:"value5,range=[1:5]"`
		Value8  float64 `key:"value8,range=[1:5],string"`
		Value9  float64 `key:"value9,range=[1:5],string"`
		Value10 float64 `key:"value10,range=[1:5],string"`
		Value11 float64 `key:"value11,range=[1:5],string"`
	}
	m := map[string]any{
		"value2":  float32(1),
		"value3":  float32(2),
		"value4":  float64(4),
		"value5":  float64(5),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(float32(1), in.Value2)
		ast.Equal(float32(2), in.Value3)
		ast.Equal(float64(4), in.Value4)
		ast.Equal(float64(5), in.Value5)
		ast.Equal(float64(1), in.Value8)
		ast.Equal(float64(2), in.Value9)
		ast.Equal(float64(4), in.Value10)
		ast.Equal(float64(5), in.Value11)
	}
}

func TestUnmarshalNumberRangeFloatLeftExclude(t *testing.T) {
	type inner struct {
		Value3  float64 `key:"value3,range=(1:5]"`
		Value4  float64 `key:"value4,range=(1:5]"`
		Value5  float64 `key:"value5,range=(1:5]"`
		Value9  float64 `key:"value9,range=(1:5],string"`
		Value10 float64 `key:"value10,range=(1:5],string"`
		Value11 float64 `key:"value11,range=(1:5],string"`
	}
	m := map[string]any{
		"value3":  float64(2),
		"value4":  float64(4),
		"value5":  float64(5),
		"value9":  "2",
		"value10": "4",
		"value11": "5",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(float64(2), in.Value3)
		ast.Equal(float64(4), in.Value4)
		ast.Equal(float64(5), in.Value5)
		ast.Equal(float64(2), in.Value9)
		ast.Equal(float64(4), in.Value10)
		ast.Equal(float64(5), in.Value11)
	}
}

func TestUnmarshalNumberRangeFloatRightExclude(t *testing.T) {
	type inner struct {
		Value2  float64 `key:"value2,range=[1:5)"`
		Value3  float64 `key:"value3,range=[1:5)"`
		Value4  float64 `key:"value4,range=[1:5)"`
		Value8  float64 `key:"value8,range=[1:5),string"`
		Value9  float64 `key:"value9,range=[1:5),string"`
		Value10 float64 `key:"value10,range=[1:5),string"`
	}
	m := map[string]any{
		"value2":  float64(1),
		"value3":  float64(2),
		"value4":  float64(4),
		"value8":  "1",
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(float64(1), in.Value2)
		ast.Equal(float64(2), in.Value3)
		ast.Equal(float64(4), in.Value4)
		ast.Equal(float64(1), in.Value8)
		ast.Equal(float64(2), in.Value9)
		ast.Equal(float64(4), in.Value10)
	}
}

func TestUnmarshalNumberRangeFloatExclude(t *testing.T) {
	type inner struct {
		Value3  float64 `key:"value3,range=(1:5)"`
		Value4  float64 `key:"value4,range=(1:5)"`
		Value9  float64 `key:"value9,range=(1:5),string"`
		Value10 float64 `key:"value10,range=(1:5),string"`
	}
	m := map[string]any{
		"value3":  float64(2),
		"value4":  float64(4),
		"value9":  "2",
		"value10": "4",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(float64(2), in.Value3)
		ast.Equal(float64(4), in.Value4)
		ast.Equal(float64(2), in.Value9)
		ast.Equal(float64(4), in.Value10)
	}
}

func TestUnmarshalNumberRangeFloatOutOfRange(t *testing.T) {
	type inner1 struct {
		Value float64 `key:"value,range=(1:5)"`
	}

	var in1 inner1
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(1),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(0),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(5),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": json.Number("6"),
	}, &in1))

	type inner2 struct {
		Value float64 `key:"value,range=[1:5)"`
	}

	var in2 inner2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(0),
	}, &in2))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(5),
	}, &in2))

	type inner3 struct {
		Value float64 `key:"value,range=(1:5]"`
	}

	var in3 inner3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(1),
	}, &in3))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(6),
	}, &in3))

	type inner4 struct {
		Value float64 `key:"value,range=[1:5]"`
	}

	var in4 inner4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(0),
	}, &in4))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"value": float64(6),
	}, &in4))
}

func TestUnmarshalRangeError(t *testing.T) {
	type inner1 struct {
		Value int `key:",range="`
	}
	var in1 inner1
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"Value": 1,
	}, &in1))

	type inner2 struct {
		Value int `key:",range=["`
	}
	var in2 inner2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"Value": 1,
	}, &in2))

	type inner3 struct {
		Value int `key:",range=[:"`
	}
	var in3 inner3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"Value": 1,
	}, &in3))

	type inner4 struct {
		Value int `key:",range=[:]"`
	}
	var in4 inner4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"Value": 1,
	}, &in4))

	type inner5 struct {
		Value int `key:",range={:]"`
	}
	var in5 inner5
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"Value": 1,
	}, &in5))

	type inner6 struct {
		Value int `key:",range=[:}"`
	}
	var in6 inner6
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"Value": 1,
	}, &in6))

	type inner7 struct {
		Value int `key:",range=[]"`
	}
	var in7 inner7
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"Value": 1,
	}, &in7))

	type inner8 struct {
		Value int `key:",range=[a:]"`
	}
	var in8 inner8
	assert.Error(t, UnmarshalKey(map[string]any{
		"Value": 1,
	}, &in8))

	type inner9 struct {
		Value int `key:",range=[:a]"`
	}
	var in9 inner9
	assert.Error(t, UnmarshalKey(map[string]any{
		"Value": 1,
	}, &in9))

	type inner10 struct {
		Value int `key:",range"`
	}
	var in10 inner10
	assert.Error(t, UnmarshalKey(map[string]any{
		"Value": 1,
	}, &in10))

	type inner11 struct {
		Value int `key:",range=[1,2]"`
	}
	var in11 inner11
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]any{
		"Value": "a",
	}, &in11))
}

func TestUnmarshalNestedMap(t *testing.T) {
	t.Run("nested map", func(t *testing.T) {
		var c struct {
			Anything map[string]map[string]string `json:"anything"`
		}
		m := map[string]any{
			"anything": map[string]map[string]any{
				"inner": {
					"id":   "1",
					"name": "any",
				},
			},
		}

		if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
			assert.Equal(t, "1", c.Anything["inner"]["id"])
		}
	})

	t.Run("nested map with slice element", func(t *testing.T) {
		var c struct {
			Anything map[string][]string `json:"anything"`
		}
		m := map[string]any{
			"anything": map[string][]any{
				"inner": {
					"id",
					"name",
				},
			},
		}

		if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
			assert.Equal(t, []string{"id", "name"}, c.Anything["inner"])
		}
	})

	t.Run("nested map with slice element error", func(t *testing.T) {
		var c struct {
			Anything map[string][]string `json:"anything"`
		}
		m := map[string]any{
			"anything": map[string][]any{
				"inner": {
					"id",
					1,
				},
			},
		}

		assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c))
	})
}

func TestUnmarshalNestedMapMismatch(t *testing.T) {
	var c struct {
		Anything map[string]map[string]map[string]string `json:"anything"`
	}
	m := map[string]any{
		"anything": map[string]map[string]any{
			"inner": {
				"name": "any",
			},
		},
	}

	assert.Error(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c))
}

func TestUnmarshalNestedMapSimple(t *testing.T) {
	var c struct {
		Anything map[string]string `json:"anything"`
	}
	m := map[string]any{
		"anything": map[string]any{
			"id":   "1",
			"name": "any",
		},
	}

	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
		assert.Equal(t, "1", c.Anything["id"])
	}
}

func TestUnmarshalNestedMapSimpleTypeMatch(t *testing.T) {
	var c struct {
		Anything map[string]string `json:"anything"`
	}
	m := map[string]any{
		"anything": map[string]string{
			"id":   "1",
			"name": "any",
		},
	}

	if assert.NoError(t, NewUnmarshaler(jsonTagKey).Unmarshal(m, &c)) {
		assert.Equal(t, "1", c.Anything["id"])
	}
}

func TestUnmarshalInheritPrimitiveUseParent(t *testing.T) {
	type (
		component struct {
			Name      string `key:"name"`
			Discovery string `key:"discovery,inherit"`
		}
		server struct {
			Discovery string    `key:"discovery"`
			Component component `key:"component"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": "localhost:8080",
		"component": map[string]any{
			"name": "test",
		},
	}, &s)) {
		assert.Equal(t, "localhost:8080", s.Discovery)
		assert.Equal(t, "localhost:8080", s.Component.Discovery)
	}
}

func TestUnmarshalInheritPrimitiveUseSelf(t *testing.T) {
	type (
		component struct {
			Name      string `key:"name"`
			Discovery string `key:"discovery,inherit"`
		}
		server struct {
			Discovery string    `key:"discovery"`
			Component component `key:"component"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": "localhost:8080",
		"component": map[string]any{
			"name":      "test",
			"discovery": "localhost:8888",
		},
	}, &s)) {
		assert.Equal(t, "localhost:8080", s.Discovery)
		assert.Equal(t, "localhost:8888", s.Component.Discovery)
	}
}

func TestUnmarshalInheritPrimitiveNotExist(t *testing.T) {
	type (
		component struct {
			Name      string `key:"name"`
			Discovery string `key:"discovery,inherit"`
		}
		server struct {
			Component component `key:"component"`
		}
	)

	var s server
	assert.Error(t, UnmarshalKey(map[string]any{
		"component": map[string]any{
			"name": "test",
		},
	}, &s))
}

func TestUnmarshalInheritStructUseParent(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Name      string    `key:"name"`
			Discovery discovery `key:"discovery,inherit"`
		}
		server struct {
			Discovery discovery `key:"discovery"`
			Component component `key:"component"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": map[string]any{
			"host": "localhost",
			"port": 8080,
		},
		"component": map[string]any{
			"name": "test",
		},
	}, &s)) {
		assert.Equal(t, "localhost", s.Discovery.Host)
		assert.Equal(t, 8080, s.Discovery.Port)
		assert.Equal(t, "localhost", s.Component.Discovery.Host)
		assert.Equal(t, 8080, s.Component.Discovery.Port)
	}
}

func TestUnmarshalInheritStructUseSelf(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Name      string    `key:"name"`
			Discovery discovery `key:"discovery,inherit"`
		}
		server struct {
			Discovery discovery `key:"discovery"`
			Component component `key:"component"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": map[string]any{
			"host": "localhost",
			"port": 8080,
		},
		"component": map[string]any{
			"name": "test",
			"discovery": map[string]any{
				"host": "remotehost",
				"port": 8888,
			},
		},
	}, &s)) {
		assert.Equal(t, "localhost", s.Discovery.Host)
		assert.Equal(t, 8080, s.Discovery.Port)
		assert.Equal(t, "remotehost", s.Component.Discovery.Host)
		assert.Equal(t, 8888, s.Component.Discovery.Port)
	}
}

func TestUnmarshalInheritStructNotExist(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Name      string    `key:"name"`
			Discovery discovery `key:"discovery,inherit"`
		}
		server struct {
			Component component `key:"component"`
		}
	)

	var s server
	assert.Error(t, UnmarshalKey(map[string]any{
		"component": map[string]any{
			"name": "test",
		},
	}, &s))
}

func TestUnmarshalInheritStructUsePartial(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Name      string    `key:"name"`
			Discovery discovery `key:"discovery,inherit"`
		}
		server struct {
			Discovery discovery `key:"discovery"`
			Component component `key:"component"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": map[string]any{
			"host": "localhost",
			"port": 8080,
		},
		"component": map[string]any{
			"name": "test",
			"discovery": map[string]any{
				"port": 8888,
			},
		},
	}, &s)) {
		assert.Equal(t, "localhost", s.Discovery.Host)
		assert.Equal(t, 8080, s.Discovery.Port)
		assert.Equal(t, "localhost", s.Component.Discovery.Host)
		assert.Equal(t, 8888, s.Component.Discovery.Port)
	}
}

func TestUnmarshalInheritStructUseSelfIncorrectType(t *testing.T) {
	type (
		discovery struct {
			Host string `key:"host"`
			Port int    `key:"port"`
		}
		component struct {
			Name      string    `key:"name"`
			Discovery discovery `key:"discovery,inherit"`
		}
		server struct {
			Discovery discovery `key:"discovery"`
			Component component `key:"component"`
		}
	)

	var s server
	assert.Error(t, UnmarshalKey(map[string]any{
		"discovery": map[string]any{
			"host": "localhost",
		},
		"component": map[string]any{
			"name": "test",
			"discovery": map[string]string{
				"host": "remotehost",
			},
		},
	}, &s))
}

func TestUnmarshaler_InheritFromGrandparent(t *testing.T) {
	type (
		component struct {
			Name      string `key:"name"`
			Discovery string `key:"discovery,inherit"`
		}
		middle struct {
			Value component `key:"value"`
		}
		server struct {
			Discovery string `key:"discovery"`
			Middle    middle `key:"middle"`
		}
	)

	var s server
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"discovery": "localhost:8080",
		"middle": map[string]any{
			"value": map[string]any{
				"name": "test",
			},
		},
	}, &s)) {
		assert.Equal(t, "localhost:8080", s.Discovery)
		assert.Equal(t, "localhost:8080", s.Middle.Value.Discovery)
	}
}

func TestUnmarshaler_InheritSequence(t *testing.T) {
	var testConf = []byte(`
Nacos:
  NamespaceId: "123"
RpcConf:
  Nacos:
    NamespaceId: "456"
  Name: hello
`)

	type (
		NacosConf struct {
			NamespaceId string
		}

		RpcConf struct {
			Nacos NacosConf `json:",inherit"`
			Name  string
		}

		Config1 struct {
			RpcConf RpcConf
			Nacos   NacosConf
		}

		Config2 struct {
			RpcConf RpcConf
			Nacos   NacosConf
		}
	)

	var c1 Config1
	if assert.NoError(t, UnmarshalYamlBytes(testConf, &c1)) {
		assert.Equal(t, "123", c1.Nacos.NamespaceId)
		assert.Equal(t, "456", c1.RpcConf.Nacos.NamespaceId)
	}

	var c2 Config2
	if assert.NoError(t, UnmarshalYamlBytes(testConf, &c2)) {
		assert.Equal(t, "123", c1.Nacos.NamespaceId)
		assert.Equal(t, "456", c1.RpcConf.Nacos.NamespaceId)
	}
}

func TestUnmarshaler_InheritNested(t *testing.T) {
	var testConf = []byte(`
Nacos:
  Value1: "123"
Server:
  Nacos:
    Value2: "456"
  Rpc:
    Nacos:
      Value3: "789"
    Name: hello
`)

	type (
		NacosConf struct {
			Value1 string `json:",optional"`
			Value2 string `json:",optional"`
			Value3 string `json:",optional"`
		}

		RpcConf struct {
			Nacos NacosConf `json:",inherit"`
			Name  string
		}

		ServerConf struct {
			Nacos NacosConf `json:",inherit"`
			Rpc   RpcConf
		}

		Config struct {
			Server ServerConf
			Nacos  NacosConf
		}
	)

	var c Config
	if assert.NoError(t, UnmarshalYamlBytes(testConf, &c)) {
		assert.Equal(t, "123", c.Nacos.Value1)
		assert.Empty(t, c.Nacos.Value2)
		assert.Empty(t, c.Nacos.Value3)
		assert.Equal(t, "123", c.Server.Nacos.Value1)
		assert.Equal(t, "456", c.Server.Nacos.Value2)
		assert.Empty(t, c.Nacos.Value3)
		assert.Equal(t, "123", c.Server.Rpc.Nacos.Value1)
		assert.Equal(t, "456", c.Server.Rpc.Nacos.Value2)
		assert.Equal(t, "789", c.Server.Rpc.Nacos.Value3)
	}
}

func TestUnmarshalValuer(t *testing.T) {
	unmarshaler := NewUnmarshaler(jsonTagKey)
	var foo string
	err := unmarshaler.UnmarshalValuer(nil, foo)
	assert.Error(t, err)
}

func TestUnmarshal_EnvString(t *testing.T) {
	t.Run("valid env", func(t *testing.T) {
		type Value struct {
			Name string `key:"name,env=TEST_NAME_STRING"`
		}

		const (
			envName = "TEST_NAME_STRING"
			envVal  = "this is a name"
		)
		t.Setenv(envName, envVal)

		var v Value
		if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
			assert.Equal(t, envVal, v.Name)
		}
	})

	t.Run("invalid env", func(t *testing.T) {
		type Value struct {
			Name string `key:"name,env=TEST_NAME_STRING=invalid"`
		}

		const (
			envName = "TEST_NAME_STRING"
			envVal  = "this is a name"
		)
		t.Setenv(envName, envVal)

		var v Value
		assert.Error(t, UnmarshalKey(emptyMap, &v))
	})
}

func TestUnmarshal_EnvStringOverwrite(t *testing.T) {
	type Value struct {
		Name string `key:"name,env=TEST_NAME_STRING"`
	}

	const (
		envName = "TEST_NAME_STRING"
		envVal  = "this is a name"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"name": "local value",
	}, &v)) {
		assert.Equal(t, envVal, v.Name)
	}
}

func TestUnmarshal_EnvInt(t *testing.T) {
	type Value struct {
		Age int `key:"age,env=TEST_NAME_INT"`
	}

	const (
		envName = "TEST_NAME_INT"
		envVal  = "123"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.Equal(t, 123, v.Age)
	}
}

func TestUnmarshal_EnvIntOverwrite(t *testing.T) {
	type Value struct {
		Age int `key:"age,env=TEST_NAME_INT"`
	}

	const (
		envName = "TEST_NAME_INT"
		envVal  = "123"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"age": 18,
	}, &v)) {
		assert.Equal(t, 123, v.Age)
	}
}

func TestUnmarshal_EnvFloat(t *testing.T) {
	type Value struct {
		Age float32 `key:"name,env=TEST_NAME_FLOAT"`
	}

	const (
		envName = "TEST_NAME_FLOAT"
		envVal  = "123.45"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.Equal(t, float32(123.45), v.Age)
	}
}

func TestUnmarshal_EnvFloatOverwrite(t *testing.T) {
	type Value struct {
		Age float32 `key:"age,env=TEST_NAME_FLOAT"`
	}

	const (
		envName = "TEST_NAME_FLOAT"
		envVal  = "123.45"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(map[string]any{
		"age": 18.5,
	}, &v)) {
		assert.Equal(t, float32(123.45), v.Age)
	}
}

func TestUnmarshal_EnvBoolTrue(t *testing.T) {
	type Value struct {
		Enable bool `key:"enable,env=TEST_NAME_BOOL_TRUE"`
	}

	const (
		envName = "TEST_NAME_BOOL_TRUE"
		envVal  = "true"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.True(t, v.Enable)
	}
}

func TestUnmarshal_EnvBoolFalse(t *testing.T) {
	type Value struct {
		Enable bool `key:"enable,env=TEST_NAME_BOOL_FALSE"`
	}

	const (
		envName = "TEST_NAME_BOOL_FALSE"
		envVal  = "false"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.False(t, v.Enable)
	}
}

func TestUnmarshal_EnvBoolBad(t *testing.T) {
	type Value struct {
		Enable bool `key:"enable,env=TEST_NAME_BOOL_BAD"`
	}

	const (
		envName = "TEST_NAME_BOOL_BAD"
		envVal  = "bad"
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestUnmarshal_EnvDuration(t *testing.T) {
	type Value struct {
		Duration time.Duration `key:"duration,env=TEST_NAME_DURATION"`
	}

	const (
		envName = "TEST_NAME_DURATION"
		envVal  = "1s"
	)
	t.Setenv(envName, envVal)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.Equal(t, time.Second, v.Duration)
	}
}

func TestUnmarshal_EnvDurationBadValue(t *testing.T) {
	type Value struct {
		Duration time.Duration `key:"duration,env=TEST_NAME_BAD_DURATION"`
	}

	const (
		envName = "TEST_NAME_BAD_DURATION"
		envVal  = "bad"
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestUnmarshal_EnvWithOptions(t *testing.T) {
	t.Run("valid options", func(t *testing.T) {
		type Value struct {
			Name string `key:"name,env=TEST_NAME_ENV_OPTIONS_MATCH,options=[abc,123,xyz]"`
		}

		const (
			envName = "TEST_NAME_ENV_OPTIONS_MATCH"
			envVal  = "123"
		)
		t.Setenv(envName, envVal)

		var v Value
		if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
			assert.Equal(t, envVal, v.Name)
		}
	})
}

func TestUnmarshal_EnvWithOptionsWrongValueBool(t *testing.T) {
	type Value struct {
		Enable bool `key:"enable,env=TEST_NAME_ENV_OPTIONS_BOOL,options=[true]"`
	}

	const (
		envName = "TEST_NAME_ENV_OPTIONS_BOOL"
		envVal  = "false"
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestUnmarshal_EnvWithOptionsWrongValueDuration(t *testing.T) {
	type Value struct {
		Duration time.Duration `key:"duration,env=TEST_NAME_ENV_OPTIONS_DURATION,options=[1s,2s,3s]"`
	}

	const (
		envName = "TEST_NAME_ENV_OPTIONS_DURATION"
		envVal  = "4s"
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestUnmarshal_EnvWithOptionsWrongValueNumber(t *testing.T) {
	type Value struct {
		Age int `key:"age,env=TEST_NAME_ENV_OPTIONS_AGE,options=[18,19,20]"`
	}

	const (
		envName = "TEST_NAME_ENV_OPTIONS_AGE"
		envVal  = "30"
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestUnmarshal_EnvWithOptionsWrongValueString(t *testing.T) {
	type Value struct {
		Name string `key:"name,env=TEST_NAME_ENV_OPTIONS_STRING,options=[abc,123,xyz]"`
	}

	const (
		envName = "TEST_NAME_ENV_OPTIONS_STRING"
		envVal  = "this is a name"
	)
	t.Setenv(envName, envVal)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestUnmarshalJsonReaderMultiArray(t *testing.T) {
	t.Run("reader multi array", func(t *testing.T) {
		type testRes struct {
			A string     `json:"a"`
			B [][]string `json:"b"`
			C []byte     `json:"c"`
		}

		var res testRes
		marshal := testRes{
			A: "133",
			B: [][]string{
				{"add", "cccd"},
				{"eeee"},
			},
			C: []byte("11122344wsss"),
		}
		bytes, err := jsonx.Marshal(marshal)
		assert.NoError(t, err)
		payload := string(bytes)
		reader := strings.NewReader(payload)
		if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
			assert.Equal(t, 2, len(res.B))
			assert.Equal(t, string(marshal.C), string(res.C))
		}
	})

	t.Run("reader multi array with error", func(t *testing.T) {
		var res struct {
			A string     `json:"a"`
			B [][]string `json:"b"`
		}
		payload := `{"a": "133", "b": ["eeee"]}`
		reader := strings.NewReader(payload)
		assert.Error(t, UnmarshalJsonReader(reader, &res))
	})
}

func TestUnmarshalJsonReaderPtrMultiArrayString(t *testing.T) {
	var res struct {
		A string      `json:"a"`
		B [][]*string `json:"b"`
	}
	payload := `{"a": "133", "b": [["add", "cccd"], ["eeee"]]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 2, len(res.B))
		assert.Equal(t, 2, len(res.B[0]))
	}
}

func TestUnmarshalJsonReaderPtrMultiArrayString_Int(t *testing.T) {
	var res struct {
		A string      `json:"a"`
		B [][]*string `json:"b"`
	}
	payload := `{"a": "133", "b": [[11, 22], [33]]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 2, len(res.B))
		assert.Equal(t, 2, len(res.B[0]))
	}
}

func TestUnmarshalJsonReaderPtrMultiArrayInt(t *testing.T) {
	var res struct {
		A string   `json:"a"`
		B [][]*int `json:"b"`
	}
	payload := `{"a": "133", "b": [[11, 22], [33]]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 2, len(res.B))
		assert.Equal(t, 2, len(res.B[0]))
	}
}

func TestUnmarshalJsonReaderPtrArray(t *testing.T) {
	var res struct {
		A string    `json:"a"`
		B []*string `json:"b"`
	}
	payload := `{"a": "133", "b": ["add", "cccd", "eeee"]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 3, len(res.B))
	}
}

func TestUnmarshalJsonReaderPtrArray_Int(t *testing.T) {
	var res struct {
		A string    `json:"a"`
		B []*string `json:"b"`
	}
	payload := `{"a": "133", "b": [11, 22, 33]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 3, len(res.B))
	}
}

func TestUnmarshalJsonReaderPtrInt(t *testing.T) {
	var res struct {
		A string    `json:"a"`
		B []*string `json:"b"`
	}
	payload := `{"a": "133", "b": [11, 22, 33]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 3, len(res.B))
	}
}

func TestUnmarshalJsonWithoutKey(t *testing.T) {
	var res struct {
		A string `json:""`
		B string `json:","`
	}
	payload := `{"A": "1", "B": "2"}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, "1", res.A)
		assert.Equal(t, "2", res.B)
	}
}

func TestUnmarshalJsonUintNegative(t *testing.T) {
	var res struct {
		A uint `json:"a"`
	}
	payload := `{"a": -1}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestUnmarshalJsonDefinedInt(t *testing.T) {
	type Int int
	var res struct {
		A Int `json:"a"`
	}
	payload := `{"a": -1}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, Int(-1), res.A)
	}
}

func TestUnmarshalJsonDefinedString(t *testing.T) {
	type String string
	var res struct {
		A String `json:"a"`
	}
	payload := `{"a": "foo"}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, String("foo"), res.A)
	}
}

func TestUnmarshalJsonDefinedStringPtr(t *testing.T) {
	type String string
	var res struct {
		A *String `json:"a"`
	}
	payload := `{"a": "foo"}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, String("foo"), *res.A)
	}
}

func TestUnmarshalJsonReaderComplex(t *testing.T) {
	type (
		MyInt      int
		MyTxt      string
		MyTxtArray []string

		Req struct {
			MyInt      MyInt      `json:"my_int"` // int.. ok
			MyTxtArray MyTxtArray `json:"my_txt_array"`
			MyTxt      MyTxt      `json:"my_txt"` // but string is not assignable
			Int        int        `json:"int"`
			Txt        string     `json:"txt"`
		}
	)
	body := `{
  "my_int": 100,
  "my_txt_array": [
    "a",
    "b"
  ],
  "my_txt": "my_txt",
  "int": 200,
  "txt": "txt"
}`
	var req Req
	if assert.NoError(t, UnmarshalJsonReader(strings.NewReader(body), &req)) {
		assert.Equal(t, MyInt(100), req.MyInt)
		assert.Equal(t, MyTxt("my_txt"), req.MyTxt)
		assert.EqualValues(t, MyTxtArray([]string{"a", "b"}), req.MyTxtArray)
		assert.Equal(t, 200, req.Int)
		assert.Equal(t, "txt", req.Txt)
	}
}

func TestUnmarshalJsonReaderArrayBool(t *testing.T) {
	var res struct {
		ID []string `json:"id"`
	}
	payload := `{"id": false}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestUnmarshalJsonReaderArrayInt(t *testing.T) {
	var res struct {
		ID []string `json:"id"`
	}
	payload := `{"id": 123}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestUnmarshalJsonReaderArrayString(t *testing.T) {
	var res struct {
		ID []string `json:"id"`
	}
	payload := `{"id": "123"}`
	reader := strings.NewReader(payload)
	assert.Error(t, UnmarshalJsonReader(reader, &res))
}

func TestGoogleUUID(t *testing.T) {
	var val struct {
		Uid    uuid.UUID    `json:"uid,optional"`
		Uidp   *uuid.UUID   `json:"uidp,optional"`
		Uidpp  **uuid.UUID  `json:"uidpp,optional"`
		Uidppp ***uuid.UUID `json:"uidppp,optional"`
	}

	t.Run("bytes", func(t *testing.T) {
		if assert.NoError(t, UnmarshalJsonBytes([]byte(`{
			"uid": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			"uidp": "a0b3d4af-4232-4c7d-b722-7ae879620518",
			"uidpp": "a0b3d4af-4232-4c7d-b722-7ae879620519",
			"uidppp": "6ba7b810-9dad-11d1-80b4-00c04fd430c9"}`), &val)) {
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", val.Uid.String())
			assert.Equal(t, "a0b3d4af-4232-4c7d-b722-7ae879620518", val.Uidp.String())
			assert.Equal(t, "a0b3d4af-4232-4c7d-b722-7ae879620519", (*val.Uidpp).String())
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c9", (**val.Uidppp).String())
		}
	})

	t.Run("map", func(t *testing.T) {
		if assert.NoError(t, UnmarshalJsonMap(map[string]any{
			"uid":    []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c1"),
			"uidp":   []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c2"),
			"uidpp":  []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c3"),
			"uidppp": []byte("6ba7b810-9dad-11d1-80b4-00c04fd430c4"),
		}, &val)) {
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c1", val.Uid.String())
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c2", val.Uidp.String())
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c3", (*val.Uidpp).String())
			assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c4", (**val.Uidppp).String())
		}
	})
}

func TestUnmarshalJsonReaderWithTypeMismatchBool(t *testing.T) {
	var req struct {
		Params map[string]bool `json:"params"`
	}
	body := `{"params":{"a":"123"}}`
	assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
}

func TestUnmarshalJsonReaderWithTypeString(t *testing.T) {
	t.Run("string type", func(t *testing.T) {
		var req struct {
			Params map[string]string `json:"params"`
		}
		body := `{"params":{"a":"b"}}`
		if assert.NoError(t, UnmarshalJsonReader(strings.NewReader(body), &req)) {
			assert.Equal(t, "b", req.Params["a"])
		}
	})

	t.Run("string type mismatch", func(t *testing.T) {
		var req struct {
			Params map[string]string `json:"params"`
		}
		body := `{"params":{"a":{"a":123}}}`
		assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
	})

	t.Run("customized string type", func(t *testing.T) {
		type myString string

		var req struct {
			Params map[string]myString `json:"params"`
		}
		body := `{"params":{"a":"b"}}`
		assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
	})
}

func TestUnmarshalJsonReaderWithMismatchType(t *testing.T) {
	type Req struct {
		Params map[string]string `json:"params"`
	}

	var req Req
	body := `{"params":{"a":{"a":123}}}`
	assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
}

func TestUnmarshalJsonReaderWithTypeBool(t *testing.T) {
	t.Run("bool type", func(t *testing.T) {
		type Req struct {
			Params map[string]bool `json:"params"`
		}

		tests := []struct {
			name   string
			input  string
			expect bool
		}{
			{
				name:   "int",
				input:  `{"params":{"a":1}}`,
				expect: true,
			},
			{
				name:   "int",
				input:  `{"params":{"a":0}}`,
				expect: false,
			},
		}

		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				var req Req
				if assert.NoError(t, UnmarshalJsonReader(strings.NewReader(test.input), &req)) {
					assert.Equal(t, test.expect, req.Params["a"])
				}
			})
		}
	})

	t.Run("bool type mismatch", func(t *testing.T) {
		type Req struct {
			Params map[string]bool `json:"params"`
		}

		tests := []struct {
			name  string
			input string
		}{
			{
				name:  "int",
				input: `{"params":{"a":123}}`,
			},
			{
				name:  "int",
				input: `{"params":{"a":"123"}}`,
			},
		}

		for _, test := range tests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				var req Req
				assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(test.input), &req))
			})
		}
	})
}

func TestUnmarshalJsonReaderWithTypeBoolMap(t *testing.T) {
	t.Run("bool map", func(t *testing.T) {
		var req struct {
			Params map[string]bool `json:"params"`
		}
		if assert.NoError(t, UnmarshalJsonMap(map[string]any{
			"params": map[string]any{
				"a": true,
			},
		}, &req)) {
			assert.Equal(t, map[string]bool{
				"a": true,
			}, req.Params)
		}
	})

	t.Run("bool map with error", func(t *testing.T) {
		var req struct {
			Params map[string]string `json:"params"`
		}
		assert.Equal(t, errTypeMismatch, UnmarshalJsonMap(map[string]any{
			"params": map[string]any{
				"a": true,
			},
		}, &req))
	})
}

func TestUnmarshalJsonBytesSliceOfMaps(t *testing.T) {
	input := []byte(`{
  "order_id": "1234567",
  "refund_reason": {
    "reason_code": [
      123,
      234
    ],
    "desc": "not wanted",
    "show_reason": [
      {
        "123": "not enough",
        "234": "closed"
      }
    ]
  },
  "product_detail": {
    "product_id": "123",
    "sku_id": "123",
    "name": "cake",
    "actual_amount": 100
  }
}`)

	type (
		RefundReasonData struct {
			ReasonCode []int               `json:"reason_code"`
			Desc       string              `json:"desc"`
			ShowReason []map[string]string `json:"show_reason"`
		}

		ProductDetailData struct {
			ProductId    string `json:"product_id"`
			SkuId        string `json:"sku_id"`
			Name         string `json:"name"`
			ActualAmount int    `json:"actual_amount"`
		}

		OrderApplyRefundReq struct {
			OrderId       string            `json:"order_id"`
			RefundReason  RefundReasonData  `json:"refund_reason,optional"`
			ProductDetail ProductDetailData `json:"product_detail,optional"`
		}
	)

	var req OrderApplyRefundReq
	assert.NoError(t, UnmarshalJsonBytes(input, &req))
}

func TestUnmarshalJsonBytesWithAnonymousField(t *testing.T) {
	type (
		Int int

		InnerConf struct {
			Name string
		}

		Conf struct {
			Int
			InnerConf
		}
	)

	var (
		input = []byte(`{"Name": "hello", "Int": 3}`)
		c     Conf
	)
	if assert.NoError(t, UnmarshalJsonBytes(input, &c)) {
		assert.Equal(t, "hello", c.Name)
		assert.Equal(t, Int(3), c.Int)
	}
}

func TestUnmarshalJsonBytesWithAnonymousFieldOptional(t *testing.T) {
	type (
		Int int

		InnerConf struct {
			Name string
		}

		Conf struct {
			Int `json:",optional"`
			InnerConf
		}
	)

	var (
		input = []byte(`{"Name": "hello", "Int": 3}`)
		c     Conf
	)
	if assert.NoError(t, UnmarshalJsonBytes(input, &c)) {
		assert.Equal(t, "hello", c.Name)
		assert.Equal(t, Int(3), c.Int)
	}
}

func TestUnmarshalJsonBytesWithAnonymousFieldBadTag(t *testing.T) {
	type (
		Int int

		InnerConf struct {
			Name string
		}

		Conf struct {
			Int `json:",optional=123"`
			InnerConf
		}
	)

	var (
		input = []byte(`{"Name": "hello", "Int": 3}`)
		c     Conf
	)
	assert.Error(t, UnmarshalJsonBytes(input, &c))
}

func TestUnmarshalJsonBytesWithAnonymousFieldBadValue(t *testing.T) {
	type (
		Int int

		InnerConf struct {
			Name string
		}

		Conf struct {
			Int
			InnerConf
		}
	)

	var (
		input = []byte(`{"Name": "hello", "Int": "3"}`)
		c     Conf
	)
	assert.Error(t, UnmarshalJsonBytes(input, &c))
}

func TestUnmarshalJsonBytesWithAnonymousFieldBadTagInStruct(t *testing.T) {
	type (
		InnerConf struct {
			Name string `json:",optional=123"`
		}

		Conf struct {
			InnerConf `json:",optional"`
		}
	)

	var (
		input = []byte(`{"Name": "hello"}`)
		c     Conf
	)
	assert.Error(t, UnmarshalJsonBytes(input, &c))
}

func TestUnmarshalJsonBytesWithAnonymousFieldNotInOptions(t *testing.T) {
	type (
		InnerConf struct {
			Name string `json:",options=[a,b]"`
		}

		Conf struct {
			InnerConf `json:",optional"`
		}
	)

	var (
		input = []byte(`{"Name": "hello"}`)
		c     Conf
	)
	assert.Error(t, UnmarshalJsonBytes(input, &c))
}

func TestUnmarshalNestedPtr(t *testing.T) {
	type inner struct {
		Int **int `key:"int"`
	}
	m := map[string]any{
		"int": 1,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.NotNil(t, in.Int)
		assert.Equal(t, 1, **in.Int)
	}
}

func TestUnmarshalStructPtrOfPtr(t *testing.T) {
	type inner struct {
		Int int `key:"int"`
	}
	m := map[string]any{
		"int": 1,
	}

	in := new(inner)
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 1, in.Int)
	}
}

func TestUnmarshalOnlyPublicVariables(t *testing.T) {
	type demo struct {
		age  int    `key:"age"`
		Name string `key:"name"`
	}

	m := map[string]any{
		"age":  3,
		"name": "go-zero",
	}

	var in demo
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 0, in.age)
		assert.Equal(t, "go-zero", in.Name)
	}
}

func TestFillDefaultUnmarshal(t *testing.T) {
	fillDefaultUnmarshal := NewUnmarshaler(jsonTagKey, WithDefault())
	t.Run("nil", func(t *testing.T) {
		type St struct{}
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, St{})
		assert.Error(t, err)
	})

	t.Run("not nil", func(t *testing.T) {
		type St struct{}
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &St{})
		assert.NoError(t, err)
	})

	t.Run("default", func(t *testing.T) {
		type St struct {
			A string `json:",default=a"`
			B string
		}
		var st St
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &st)
		assert.NoError(t, err)
		assert.Equal(t, "a", st.A)
	})

	t.Run("env", func(t *testing.T) {
		type St struct {
			A string `json:",default=a"`
			B string
			C string `json:",env=TEST_C"`
		}
		t.Setenv("TEST_C", "c")

		var st St
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &st)
		assert.NoError(t, err)
		assert.Equal(t, "a", st.A)
		assert.Equal(t, "c", st.C)
	})

	t.Run("optional !", func(t *testing.T) {
		var st struct {
			A string `json:",optional"`
			B string `json:",optional=!A"`
		}
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &st)
		assert.NoError(t, err)
	})

	t.Run("has value", func(t *testing.T) {
		type St struct {
			A string `json:",default=a"`
			B string
		}
		var st = St{
			A: "b",
		}
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &st)
		assert.Error(t, err)
	})

	t.Run("handling struct", func(t *testing.T) {
		type St struct {
			A string `json:",default=a"`
			B string
		}
		type St2 struct {
			St
			St1   St
			St3   *St
			C     string `json:",default=c"`
			D     string
			Child *St2
		}
		var st2 St2
		err := fillDefaultUnmarshal.Unmarshal(map[string]any{}, &st2)
		assert.NoError(t, err)
		assert.Equal(t, "a", st2.St.A)
		assert.Equal(t, "a", st2.St1.A)
		assert.Nil(t, st2.St3)
		assert.Equal(t, "c", st2.C)
		assert.Nil(t, st2.Child)
	})
}

func Test_UnmarshalMap(t *testing.T) {
	t.Run("type mismatch", func(t *testing.T) {
		type Customer struct {
			Names map[int]string `key:"names"`
		}

		input := map[string]any{
			"names": map[string]any{
				"19": "Tom",
			},
		}

		var customer Customer
		assert.ErrorIs(t, UnmarshalKey(input, &customer), errTypeMismatch)
	})

	t.Run("map type mismatch", func(t *testing.T) {
		type Customer struct {
			Names struct {
				Values map[string]string
			} `key:"names"`
		}

		input := map[string]any{
			"names": map[string]string{
				"19": "Tom",
			},
		}

		var customer Customer
		assert.ErrorIs(t, UnmarshalKey(input, &customer), errTypeMismatch)
	})

	t.Run("map from string", func(t *testing.T) {
		type Customer struct {
			Names map[string]string `key:"names,string"`
		}

		input := map[string]any{
			"names": `{"name": "Tom"}`,
		}

		var customer Customer
		assert.NoError(t, UnmarshalKey(input, &customer))
		assert.Equal(t, "Tom", customer.Names["name"])
	})

	t.Run("map from string with error", func(t *testing.T) {
		type Customer struct {
			Names map[string]any `key:"names,string"`
		}

		input := map[string]any{
			"names": `"name"`,
		}

		var customer Customer
		assert.Error(t, UnmarshalKey(input, &customer))
	})
}

func TestUnmarshaler_Unmarshal(t *testing.T) {
	t.Run("not struct", func(t *testing.T) {
		var i int
		unmarshaler := NewUnmarshaler(jsonTagKey)
		err := unmarshaler.UnmarshalValuer(nil, &i)
		assert.Error(t, err)
	})

	t.Run("slice element missing error", func(t *testing.T) {
		type inner struct {
			S []struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			} `json:"s"`
		}
		content := []byte(`{"s": [{"name": "foo"}]}`)
		var s inner
		err := UnmarshalJsonBytes(content, &s)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s[0].age")
	})

	t.Run("map element missing error", func(t *testing.T) {
		type inner struct {
			S map[string]struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			} `json:"s"`
		}
		content := []byte(`{"s": {"a":{"name": "foo"}}}`)
		var s inner
		err := UnmarshalJsonBytes(content, &s)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "s[a].age")
	})
}

// TestUnmarshalerProcessFieldPrimitiveWithJSONNumber test the number type check.
func TestUnmarshalerProcessFieldPrimitiveWithJSONNumber(t *testing.T) {
	t.Run("wrong type", func(t *testing.T) {
		expectValue := "1"
		realValue := 1
		fieldType := reflect.TypeOf(expectValue)
		value := reflect.ValueOf(&realValue) // pass a pointer to the value
		v := json.Number(expectValue)
		m := NewUnmarshaler("field")
		err := m.processFieldPrimitiveWithJSONNumber(fieldType, value.Elem(), v,
			&fieldOptionsWithContext{}, "field")
		assert.Error(t, err)
		assert.Equal(t, `type mismatch for field "field", expect "string", actual "number"`, err.Error())
	})

	t.Run("right type", func(t *testing.T) {
		expectValue := int64(1)
		realValue := int64(1)
		fieldType := reflect.TypeOf(expectValue)
		value := reflect.ValueOf(&realValue) // pass a pointer to the value
		v := json.Number(strconv.FormatInt(expectValue, 10))
		m := NewUnmarshaler("field")
		err := m.processFieldPrimitiveWithJSONNumber(fieldType, value.Elem(), v,
			&fieldOptionsWithContext{}, "field")
		assert.NoError(t, err)
	})
}

func TestGetValueWithChainedKeys(t *testing.T) {
	t.Run("no key", func(t *testing.T) {
		_, ok := getValueWithChainedKeys(nil, []string{})
		assert.False(t, ok)
	})

	t.Run("one key", func(t *testing.T) {
		v, ok := getValueWithChainedKeys(mockValuerWithParent{
			value: "bar",
			ok:    true,
		}, []string{"foo"})
		assert.True(t, ok)
		assert.Equal(t, "bar", v)
	})

	t.Run("two keys", func(t *testing.T) {
		v, ok := getValueWithChainedKeys(mockValuerWithParent{
			value: map[string]any{
				"bar": "baz",
			},
			ok: true,
		}, []string{"foo", "bar"})
		assert.True(t, ok)
		assert.Equal(t, "baz", v)
	})

	t.Run("two keys not found", func(t *testing.T) {
		_, ok := getValueWithChainedKeys(mockValuerWithParent{
			value: "bar",
			ok:    false,
		}, []string{"foo", "bar"})
		assert.False(t, ok)
	})

	t.Run("two keys type mismatch", func(t *testing.T) {
		_, ok := getValueWithChainedKeys(mockValuerWithParent{
			value: "bar",
			ok:    true,
		}, []string{"foo", "bar"})
		assert.False(t, ok)
	})
}

func TestUnmarshalFromStringSliceForTypeMismatch(t *testing.T) {
	var v struct {
		Values map[string][]string `key:"values"`
	}
	assert.Error(t, UnmarshalKey(map[string]any{
		"values": map[string]any{
			"foo": "bar",
		},
	}, &v))
}

func TestUnmarshalWithFromArray(t *testing.T) {
	t.Run("array", func(t *testing.T) {
		var v struct {
			Value []string `key:"value"`
		}
		unmarshaler := NewUnmarshaler("key", WithFromArray())
		if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
			"value": []string{"foo", "bar"},
		}, &v)) {
			assert.ElementsMatch(t, []string{"foo", "bar"}, v.Value)
		}
	})

	t.Run("not array", func(t *testing.T) {
		var v struct {
			Value string `key:"value"`
		}
		unmarshaler := NewUnmarshaler("key", WithFromArray())
		if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
			"value": []string{"foo"},
		}, &v)) {
			assert.Equal(t, "foo", v.Value)
		}
	})

	t.Run("not array and empty", func(t *testing.T) {
		var v struct {
			Value string `key:"value"`
		}
		unmarshaler := NewUnmarshaler("key", WithFromArray())
		if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
			"value": []string{""},
		}, &v)) {
			assert.Empty(t, v.Value)
		}
	})

	t.Run("not array and no value", func(t *testing.T) {
		var v struct {
			Value string `key:"value"`
		}
		unmarshaler := NewUnmarshaler("key", WithFromArray())
		assert.Error(t, unmarshaler.Unmarshal(map[string]any{}, &v))
	})

	t.Run("not array and no value and optional", func(t *testing.T) {
		var v struct {
			Value string `key:"value,optional"`
		}
		unmarshaler := NewUnmarshaler("key", WithFromArray())
		if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{}, &v)) {
			assert.Empty(t, v.Value)
		}
	})
}

func TestUnmarshalWithOpaqueKeys(t *testing.T) {
	var v struct {
		Opaque string `key:"opaque.key"`
		Value  string `key:"value"`
	}
	unmarshaler := NewUnmarshaler("key", WithOpaqueKeys())
	if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
		"opaque.key": "foo",
		"value":      "bar",
	}, &v)) {
		assert.Equal(t, "foo", v.Opaque)
		assert.Equal(t, "bar", v.Value)
	}
}

func TestUnmarshalWithIgnoreFields(t *testing.T) {
	type (
		Foo struct {
			Value        string
			IgnoreString string `json:"-"`
			IgnoreInt    int    `json:"-"`
		}

		Bar struct {
			Foo1 Foo
			Foo2 *Foo
			Foo3 []Foo
			Foo4 []*Foo
			Foo5 map[string]Foo
			Foo6 map[string]Foo
		}

		Bar1 struct {
			Foo `json:"-"`
		}

		Bar2 struct {
			*Foo `json:"-"`
		}
	)

	var bar Bar
	unmarshaler := NewUnmarshaler(jsonTagKey)
	if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
		"Foo1": map[string]any{
			"Value":        "foo",
			"IgnoreString": "any",
			"IgnoreInt":    2,
		},
		"Foo2": map[string]any{
			"Value":        "foo",
			"IgnoreString": "any",
			"IgnoreInt":    2,
		},
		"Foo3": []map[string]any{
			{
				"Value":        "foo",
				"IgnoreString": "any",
				"IgnoreInt":    2,
			},
		},
		"Foo4": []map[string]any{
			{
				"Value":        "foo",
				"IgnoreString": "any",
				"IgnoreInt":    2,
			},
		},
		"Foo5": map[string]any{
			"key": map[string]any{
				"Value":        "foo",
				"IgnoreString": "any",
				"IgnoreInt":    2,
			},
		},
		"Foo6": map[string]any{
			"key": map[string]any{
				"Value":        "foo",
				"IgnoreString": "any",
				"IgnoreInt":    2,
			},
		},
	}, &bar)) {
		assert.Equal(t, "foo", bar.Foo1.Value)
		assert.Empty(t, bar.Foo1.IgnoreString)
		assert.Equal(t, 0, bar.Foo1.IgnoreInt)
		assert.Equal(t, "foo", bar.Foo2.Value)
		assert.Empty(t, bar.Foo2.IgnoreString)
		assert.Equal(t, 0, bar.Foo2.IgnoreInt)
		assert.Equal(t, "foo", bar.Foo3[0].Value)
		assert.Empty(t, bar.Foo3[0].IgnoreString)
		assert.Equal(t, 0, bar.Foo3[0].IgnoreInt)
		assert.Equal(t, "foo", bar.Foo4[0].Value)
		assert.Empty(t, bar.Foo4[0].IgnoreString)
		assert.Equal(t, 0, bar.Foo4[0].IgnoreInt)
		assert.Equal(t, "foo", bar.Foo5["key"].Value)
		assert.Empty(t, bar.Foo5["key"].IgnoreString)
		assert.Equal(t, 0, bar.Foo5["key"].IgnoreInt)
		assert.Equal(t, "foo", bar.Foo6["key"].Value)
		assert.Empty(t, bar.Foo6["key"].IgnoreString)
		assert.Equal(t, 0, bar.Foo6["key"].IgnoreInt)
	}

	var bar1 Bar1
	if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
		"Value":        "foo",
		"IgnoreString": "any",
		"IgnoreInt":    2,
	}, &bar1)) {
		assert.Empty(t, bar1.Value)
		assert.Empty(t, bar1.IgnoreString)
		assert.Equal(t, 0, bar1.IgnoreInt)
	}

	var bar2 Bar2
	if assert.NoError(t, unmarshaler.Unmarshal(map[string]any{
		"Value":        "foo",
		"IgnoreString": "any",
		"IgnoreInt":    2,
	}, &bar2)) {
		assert.Nil(t, bar2.Foo)
	}
}

func TestUnmarshal_Unmarshaler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		v := struct {
			Foo *mockUnmarshaler `json:"name"`
		}{}
		body := `{"name": "hello"}`
		assert.NoError(t, UnmarshalJsonBytes([]byte(body), &v))
		assert.Equal(t, "hello", v.Foo.Name)
	})

	t.Run("failure", func(t *testing.T) {
		v := struct {
			Foo *mockUnmarshalerWithError `json:"name"`
		}{}
		body := `{"name": "hello"}`
		assert.Error(t, UnmarshalJsonBytes([]byte(body), &v))
	})

	t.Run("not json unmarshaler", func(t *testing.T) {
		v := struct {
			Foo *struct {
				Name string
			} `key:"name"`
		}{}
		u := NewUnmarshaler(defaultKeyName)
		assert.Error(t, u.Unmarshal(map[string]any{
			"name": "hello",
		}, &v))
	})

	t.Run("not with json key", func(t *testing.T) {
		v := struct {
			Foo *mockUnmarshaler `json:"name"`
		}{}
		u := NewUnmarshaler(defaultKeyName)
		// with different key, ignore
		assert.NoError(t, u.Unmarshal(map[string]any{
			"name": "hello",
		}, &v))
		assert.Nil(t, v.Foo)
	})
}

func TestParseJsonStringValue(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		type GoodsInfo struct {
			Sku int64 `json:"sku,optional"`
		}

		type GetReq struct {
			GoodsList []*GoodsInfo `json:"goods_list"`
		}

		input := map[string]any{"goods_list": "[{\"sku\":11},{\"sku\":22}]"}
		var v GetReq
		assert.NotPanics(t, func() {
			assert.NoError(t, UnmarshalJsonMap(input, &v))
			assert.Equal(t, 2, len(v.GoodsList))
			assert.ElementsMatch(t, []int64{11, 22}, []int64{v.GoodsList[0].Sku, v.GoodsList[1].Sku})
		})
	})

	t.Run("string with invalid type", func(t *testing.T) {
		type GetReq struct {
			GoodsList []*int `json:"goods_list"`
		}

		input := map[string]any{"goods_list": "[{\"sku\":11},{\"sku\":22}]"}
		var v GetReq
		assert.NotPanics(t, func() {
			assert.Error(t, UnmarshalJsonMap(input, &v))
		})
	})
}

func BenchmarkDefaultValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var a struct {
			Ints []int    `json:"ints,default=[1,2,3]"`
			Strs []string `json:"strs,default=[foo,bar,baz]"`
		}
		_ = UnmarshalJsonMap(nil, &a)
		if len(a.Strs) != 3 || len(a.Ints) != 3 {
			b.Fatal("failed")
		}
	}
}

func BenchmarkUnmarshalString(b *testing.B) {
	type inner struct {
		Value string `key:"value"`
	}
	m := map[string]any{
		"value": "first",
	}

	for i := 0; i < b.N; i++ {
		var in inner
		if err := UnmarshalKey(m, &in); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalStruct(b *testing.B) {
	b.ReportAllocs()

	m := map[string]any{
		"Ids": []map[string]any{
			{
				"First":  1,
				"Second": 2,
			},
		},
	}

	for i := 0; i < b.N; i++ {
		var v struct {
			Ids []struct {
				First  int
				Second int
			}
		}
		if err := UnmarshalKey(m, &v); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapToStruct(b *testing.B) {
	data := map[string]any{
		"valid": "1",
		"age":   "5",
		"name":  "liao",
	}
	type anonymous struct {
		Valid bool
		Age   int
		Name  string
	}

	for i := 0; i < b.N; i++ {
		var an anonymous
		if valid, ok := data["valid"]; ok {
			an.Valid = valid == "1"
		}
		if age, ok := data["age"]; ok {
			ages, _ := age.(string)
			an.Age, _ = strconv.Atoi(ages)
		}
		if name, ok := data["name"]; ok {
			names, _ := name.(string)
			an.Name = names
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	data := map[string]any{
		"valid": "1",
		"age":   "5",
		"name":  "liao",
	}
	type anonymous struct {
		Valid bool   `key:"valid,string"`
		Age   int    `key:"age,string"`
		Name  string `key:"name"`
	}

	for i := 0; i < b.N; i++ {
		var an anonymous
		UnmarshalKey(data, &an)
	}
}

type mockValuerWithParent struct {
	parent valuerWithParent
	value  any
	ok     bool
}

func (m mockValuerWithParent) Value(_ string) (any, bool) {
	return m.value, m.ok
}

func (m mockValuerWithParent) Parent() valuerWithParent {
	return m.parent
}

type mockUnmarshaler struct {
	Name string
}

func (m *mockUnmarshaler) UnmarshalJSON(b []byte) error {
	m.Name = string(b)
	return nil
}

type mockUnmarshalerWithError struct {
	Name string
}

func (m *mockUnmarshalerWithError) UnmarshalJSON(b []byte) error {
	return errors.New("foo")
}
