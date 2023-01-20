package mapping

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
)

// because json.Number doesn't support strconv.ParseUint(...),
// so we only can test to 62 bits.
const maxUintBitsToTest = 62

func TestUnmarshalWithFullNameNotStruct(t *testing.T) {
	var s map[string]interface{}
	content := []byte(`{"name":"xiaoming"}`)
	err := UnmarshalJsonBytes(content, &s)
	assert.Equal(t, errTypeMismatch, err)
}

func TestUnmarshalValueNotSettable(t *testing.T) {
	var s map[string]interface{}
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
	m := map[string]interface{}{
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

func TestUnmarshalWithoutTagNameWithCanonicalKey(t *testing.T) {
	type inner struct {
		Name string `key:"name"`
	}
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	}
	m := map[string]interface{}{
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
	}
}

func TestUnmarshalIntPtr(t *testing.T) {
	type inner struct {
		Int *int `key:"int"`
	}
	m := map[string]interface{}{
		"int": 1,
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.NotNil(t, in.Int)
		assert.Equal(t, 1, *in.Int)
	}
}

func TestUnmarshalIntSliceOfPtr(t *testing.T) {
	type inner struct {
		Ints  []*int  `key:"ints"`
		Intps []**int `key:"intps"`
	}
	m := map[string]interface{}{
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
}

func TestUnmarshalIntWithDefault(t *testing.T) {
	type inner struct {
		Int   int   `key:"int,default=5"`
		Intp  *int  `key:"intp,default=5"`
		Intpp **int `key:"intpp,default=5"`
	}
	m := map[string]interface{}{
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
		m := map[string]interface{}{
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

	t.Run("int with ptr", func(t *testing.T) {
		type inner struct {
			Int *int64 `key:"int"`
		}
		m := map[string]interface{}{
			"int": json.Number("1"),
		}

		var in inner
		if assert.NoError(t, UnmarshalKey(m, &in)) {
			assert.Equal(t, int64(1), *in.Int)
		}
	})

	t.Run("int with ptr of ptr", func(t *testing.T) {
		type inner struct {
			Int **int64 `key:"int"`
		}
		m := map[string]interface{}{
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
		m := map[string]interface{}{
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
		m := map[string]interface{}{
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
		m := map[string]interface{}{
			"int": StrType("1"),
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
	assert.NotNil(t, UnmarshalKey(map[string]interface{}{}, &in))
}

func TestUnmarshalBoolSliceNil(t *testing.T) {
	type inner struct {
		Bools []bool `key:"bools,optional"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{}, &in)) {
		assert.Nil(t, in.Bools)
	}
}

func TestUnmarshalBoolSliceNilExplicit(t *testing.T) {
	type inner struct {
		Bools []bool `key:"bools,optional"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{
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
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{
		"bools": []bool{},
	}, &in)) {
		assert.Empty(t, in.Bools)
	}
}

func TestUnmarshalBoolSliceWithDefault(t *testing.T) {
	type inner struct {
		Bools []bool `key:"bools,default=[true,false]"`
	}

	var in inner
	if assert.NoError(t, UnmarshalKey(nil, &in)) {
		assert.ElementsMatch(t, []bool{true, false}, in.Bools)
	}
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
	m := map[string]interface{}{
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
		Float64      float64 `key:"float64"`
		Float64Str   float64 `key:"float64str,string"`
		DefaultFloat float32 `key:"defaultfloat,default=5.5"`
		Optional     float32 `key:",optional"`
	}
	m := map[string]interface{}{
		"float32":    float32(1.5),
		"float32str": "2.5",
		"float64":    float64(3.5),
		"float64str": "4.5",
	}

	var in inner
	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &in)) {
		ast.Equal(float32(1.5), in.Float32)
		ast.Equal(float32(2.5), in.Float32Str)
		ast.Equal(3.5, in.Float64)
		ast.Equal(4.5, in.Float64Str)
		ast.Equal(float32(5.5), in.DefaultFloat)
	}
}

func TestUnmarshalInt64Slice(t *testing.T) {
	var v struct {
		Ages  []int64 `key:"ages"`
		Slice []int64 `key:"slice"`
	}
	m := map[string]interface{}{
		"ages":  []int64{1, 2},
		"slice": []interface{}{},
	}

	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &v)) {
		ast.ElementsMatch([]int64{1, 2}, v.Ages)
		ast.Equal([]int64{}, v.Slice)
	}
}

func TestUnmarshalIntSlice(t *testing.T) {
	var v struct {
		Ages  []int `key:"ages"`
		Slice []int `key:"slice"`
	}
	m := map[string]interface{}{
		"ages":  []int{1, 2},
		"slice": []interface{}{},
	}

	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &v)) {
		ast.ElementsMatch([]int{1, 2}, v.Ages)
		ast.Equal([]int{}, v.Slice)
	}
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalStringSliceFromString(t *testing.T) {
	var v struct {
		Names []string `key:"names"`
	}
	m := map[string]interface{}{
		"names": `["first", "second"]`,
	}

	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &v)) {
		ast.Equal(2, len(v.Names))
		ast.Equal("first", v.Names[0])
		ast.Equal("second", v.Names[1])
	}
}

func TestUnmarshalIntSliceFromString(t *testing.T) {
	var v struct {
		Values []int `key:"values"`
	}
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	var v struct {
		Sort map[string]string `key:"sort"`
	}
	m := map[string]interface{}{
		"sort": CustomStringer(`"value":"ascend","emptyStr":""`),
	}

	ast := assert.New(t)
	if ast.NoError(UnmarshalKey(m, &v)) {
		ast.Equal(2, len(v.Sort))
		ast.Equal("ascend", v.Sort["value"])
		ast.Equal("", v.Sort["emptyStr"])
	}
}

func TestUnmarshalStringMapFromUnsupportedType(t *testing.T) {
	var v struct {
		Sort map[string]string `key:"sort"`
	}
	m := map[string]interface{}{
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
	m := map[string]interface{}{
		"sort":  `{"value":"ascend","emptyStr":""}`,
		"psort": `{"value":"ascend","emptyStr":""}`,
	}

	ast := assert.New(t)
	ast.Error(UnmarshalKey(m, &v))
}

func TestUnmarshalStringMapFromString(t *testing.T) {
	var v struct {
		Sort map[string]string `key:"sort"`
	}
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
		"name": "kevin",
		"address": map[string]interface{}{
			"city":    "shanghai",
			"zipcode": "200000",
		},
		"addressp": map[string]interface{}{
			"city":    "beijing",
			"zipcode": "300000",
		},
		"addresspp": map[string]interface{}{
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
			m := map[string]interface{}{
				"name": "kevin",
				"address": map[string]interface{}{
					"city": "shanghai",
				},
			}
			for k, v := range test.input {
				m["address"].(map[string]interface{})[k] = v
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
			m := map[string]interface{}{
				"name": "kevin",
				"address": map[string]interface{}{
					"city": "shanghai",
				},
			}
			for k, v := range test.input {
				m["address"].(map[string]interface{})[k] = v
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
	type address struct {
		Optional        string `key:",optional"`
		OptionalDepends string `key:",optional=!Optional"`
	}
	type inner struct {
		Name    string  `key:"name"`
		Address address `key:"address"`
	}

	m := map[string]interface{}{
		"name": "kevin",
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
}

func TestUnmarshalStructOptionalDependsNotNested(t *testing.T) {
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

	m := map[string]interface{}{
		"name": "kevin",
	}

	var in inner
	assert.Error(t, UnmarshalKey(m, &in))
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

	m := map[string]interface{}{
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

	m := map[string]interface{}{
		"name":    "kevin",
		"address": map[string]interface{}{},
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

	m := map[string]interface{}{
		"name":    "kevin",
		"address": map[string]interface{}{},
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

	m := map[string]interface{}{
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

	m := map[string]interface{}{
		"name":    "kevin",
		"address": map[string]interface{}{},
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
			m := map[string]interface{}{
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
	m := map[string]interface{}{
		"name": "kevin",
		"address": map[string]interface{}{
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
	m := map[string]interface{}{
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
		m := map[string]interface{}{
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
		m := map[string]interface{}{
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
		m := map[string]interface{}{
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
	m := map[string]interface{}{
		"Ids": map[string]bool{"first": true},
	}
	var v struct {
		Ids map[string]bool
	}
	if assert.NoError(t, UnmarshalKey(m, &v)) {
		assert.True(t, v.Ids["first"])
	}
}

func TestUnmarshalMapOfStructError(t *testing.T) {
	m := map[string]interface{}{
		"Ids": map[string]interface{}{"first": "second"},
	}
	var v struct {
		Ids map[string]struct {
			Name string
		}
	}
	assert.Error(t, UnmarshalKey(m, &v))
}

func TestUnmarshalSlice(t *testing.T) {
	m := map[string]interface{}{
		"Ids": []interface{}{"first", "second"},
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
}

func TestUnmarshalSliceOfStruct(t *testing.T) {
	m := map[string]interface{}{
		"Ids": []map[string]interface{}{
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
}

func TestUnmarshalWithStringOptionsCorrect(t *testing.T) {
	type inner struct {
		Value   string `key:"value,options=first|second"`
		Foo     string `key:"foo,options=[bar,baz]"`
		Correct string `key:"correct,options=1|2"`
	}
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{}

	var in inner
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, "second", in.Value)
	}
}

func TestUnmarshalWithOptionsAndSet(t *testing.T) {
	type inner struct {
		Value string `key:"value,options=first|second|third,default=second"`
	}
	m := map[string]interface{}{
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
	m := map[string]interface{}{
		"Persons": map[string]interface{}{
			"first": map[string]interface{}{
				"ID": 1,
			},
		},
	}

	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &c)) {
		assert.Equal(t, 1, c.ID)
	}
}

func TestUnmarhsalNestedKeyArray(t *testing.T) {
	var c struct {
		First []struct {
			ID int
		} `json:"Persons.first"`
	}
	m := map[string]interface{}{
		"Persons": map[string]interface{}{
			"first": []map[string]interface{}{
				{"ID": 1},
				{"ID": 2},
			},
		},
	}

	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &c)) {
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
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler("json").Unmarshal(m, &b))
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
	m := map[string]interface{}{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"n": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler("json").Unmarshal(m, &b))
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
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{}

	var b Bar
	assert.Error(t, NewUnmarshaler("json").Unmarshal(m, &b))
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
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler("json").Unmarshal(m, &b))
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
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler("json").Unmarshal(m, &b))
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
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{}

	var b Bar
	assert.Error(t, NewUnmarshaler("json").Unmarshal(m, &b))
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
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"n": "anything",
	}

	var b Bar
	assert.Error(t, NewUnmarshaler("json").Unmarshal(m, &b))
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
	m := map[string]interface{}{
		"n": "kevin",
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"v": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"Foo": map[string]interface{}{
			"n": "name",
			"v": "anything",
		},
	}

	var b Bar
	assert.Error(t, NewUnmarshaler("json").Unmarshal(m, &b))
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
	m := map[string]interface{}{
		"Inner": map[string]interface{}{
			"v": "anything",
		},
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"Name": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"Inner": map[string]interface{}{
			"Hosts": hosts,
			"Key":   "key",
		},
		"Name": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"Inner": map[string]interface{}{
			"Host": "thehost",
			"Key":  "thekey",
		},
		"Name": "anything",
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"foo": map[string]interface{}{
			"v": "anything",
		},
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
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
	m := map[string]interface{}{
		"Inner": map[string]interface{}{
			"v": "anything",
		},
	}

	var b Bar
	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &b)) {
		assert.Equal(t, "anything", b.Inner.Value)
	}
}

func TestUnmarshalInt2String(t *testing.T) {
	type inner struct {
		Int string `key:"int"`
	}
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m = map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(1),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(0),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(5),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": json.Number("6"),
	}, &in1))

	type inner2 struct {
		Value int64 `key:"value,optional,range=[1:5)"`
	}

	var in2 inner2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(0),
	}, &in2))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(5),
	}, &in2))

	type inner3 struct {
		Value int64 `key:"value,range=(1:5]"`
	}

	var in3 inner3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(1),
	}, &in3))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(6),
	}, &in3))

	type inner4 struct {
		Value int64 `key:"value,range=[1:5]"`
	}

	var in4 inner4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": int64(0),
	}, &in4))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
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
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(1),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(0),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(5),
	}, &in1))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": json.Number("6"),
	}, &in1))

	type inner2 struct {
		Value float64 `key:"value,range=[1:5)"`
	}

	var in2 inner2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(0),
	}, &in2))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(5),
	}, &in2))

	type inner3 struct {
		Value float64 `key:"value,range=(1:5]"`
	}

	var in3 inner3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(1),
	}, &in3))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(6),
	}, &in3))

	type inner4 struct {
		Value float64 `key:"value,range=[1:5]"`
	}

	var in4 inner4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(0),
	}, &in4))
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"value": float64(6),
	}, &in4))
}

func TestUnmarshalRangeError(t *testing.T) {
	type inner1 struct {
		Value int `key:",range="`
	}
	var in1 inner1
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in1))

	type inner2 struct {
		Value int `key:",range=["`
	}
	var in2 inner2
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in2))

	type inner3 struct {
		Value int `key:",range=[:"`
	}
	var in3 inner3
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in3))

	type inner4 struct {
		Value int `key:",range=[:]"`
	}
	var in4 inner4
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in4))

	type inner5 struct {
		Value int `key:",range={:]"`
	}
	var in5 inner5
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in5))

	type inner6 struct {
		Value int `key:",range=[:}"`
	}
	var in6 inner6
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in6))

	type inner7 struct {
		Value int `key:",range=[]"`
	}
	var in7 inner7
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in7))

	type inner8 struct {
		Value int `key:",range=[a:]"`
	}
	var in8 inner8
	assert.Error(t, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in8))

	type inner9 struct {
		Value int `key:",range=[:a]"`
	}
	var in9 inner9
	assert.Error(t, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in9))

	type inner10 struct {
		Value int `key:",range"`
	}
	var in10 inner10
	assert.Error(t, UnmarshalKey(map[string]interface{}{
		"Value": 1,
	}, &in10))

	type inner11 struct {
		Value int `key:",range=[1,2]"`
	}
	var in11 inner11
	assert.Equal(t, errNumberRange, UnmarshalKey(map[string]interface{}{
		"Value": "a",
	}, &in11))
}

func TestUnmarshalNestedMap(t *testing.T) {
	var c struct {
		Anything map[string]map[string]string `json:"anything"`
	}
	m := map[string]interface{}{
		"anything": map[string]map[string]interface{}{
			"inner": {
				"id":   "1",
				"name": "any",
			},
		},
	}

	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &c)) {
		assert.Equal(t, "1", c.Anything["inner"]["id"])
	}
}

func TestUnmarshalNestedMapMismatch(t *testing.T) {
	var c struct {
		Anything map[string]map[string]map[string]string `json:"anything"`
	}
	m := map[string]interface{}{
		"anything": map[string]map[string]interface{}{
			"inner": {
				"name": "any",
			},
		},
	}

	assert.Error(t, NewUnmarshaler("json").Unmarshal(m, &c))
}

func TestUnmarshalNestedMapSimple(t *testing.T) {
	var c struct {
		Anything map[string]string `json:"anything"`
	}
	m := map[string]interface{}{
		"anything": map[string]interface{}{
			"id":   "1",
			"name": "any",
		},
	}

	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &c)) {
		assert.Equal(t, "1", c.Anything["id"])
	}
}

func TestUnmarshalNestedMapSimpleTypeMatch(t *testing.T) {
	var c struct {
		Anything map[string]string `json:"anything"`
	}
	m := map[string]interface{}{
		"anything": map[string]string{
			"id":   "1",
			"name": "any",
		},
	}

	if assert.NoError(t, NewUnmarshaler("json").Unmarshal(m, &c)) {
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
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{
		"discovery": "localhost:8080",
		"component": map[string]interface{}{
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
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{
		"discovery": "localhost:8080",
		"component": map[string]interface{}{
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
	assert.Error(t, UnmarshalKey(map[string]interface{}{
		"component": map[string]interface{}{
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
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{
		"discovery": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
		"component": map[string]interface{}{
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
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{
		"discovery": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
		"component": map[string]interface{}{
			"name": "test",
			"discovery": map[string]interface{}{
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
	assert.Error(t, UnmarshalKey(map[string]interface{}{
		"component": map[string]interface{}{
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
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{
		"discovery": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
		"component": map[string]interface{}{
			"name": "test",
			"discovery": map[string]interface{}{
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
	assert.Error(t, UnmarshalKey(map[string]interface{}{
		"discovery": map[string]interface{}{
			"host": "localhost",
		},
		"component": map[string]interface{}{
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
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{
		"discovery": "localhost:8080",
		"middle": map[string]interface{}{
			"value": map[string]interface{}{
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
	type Value struct {
		Name string `key:"name,env=TEST_NAME_STRING"`
	}

	const (
		envName = "TEST_NAME_STRING"
		envVal  = "this is a name"
	)
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.Equal(t, envVal, v.Name)
	}
}

func TestUnmarshal_EnvStringOverwrite(t *testing.T) {
	type Value struct {
		Name string `key:"name,env=TEST_NAME_STRING"`
	}

	const (
		envName = "TEST_NAME_STRING"
		envVal  = "this is a name"
	)
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

	var v Value
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{
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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

	var v Value
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{
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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

	var v Value
	if assert.NoError(t, UnmarshalKey(map[string]interface{}{
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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestUnmarshal_EnvWithOptions(t *testing.T) {
	type Value struct {
		Name string `key:"name,env=TEST_NAME_ENV_OPTIONS_MATCH,options=[abc,123,xyz]"`
	}

	const (
		envName = "TEST_NAME_ENV_OPTIONS_MATCH"
		envVal  = "123"
	)
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

	var v Value
	if assert.NoError(t, UnmarshalKey(emptyMap, &v)) {
		assert.Equal(t, envVal, v.Name)
	}
}

func TestUnmarshal_EnvWithOptionsWrongValueBool(t *testing.T) {
	type Value struct {
		Enable bool `key:"enable,env=TEST_NAME_ENV_OPTIONS_BOOL,options=[true]"`
	}

	const (
		envName = "TEST_NAME_ENV_OPTIONS_BOOL"
		envVal  = "false"
	)
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

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
	os.Setenv(envName, envVal)
	defer os.Unsetenv(envName)

	var v Value
	assert.Error(t, UnmarshalKey(emptyMap, &v))
}

func TestUnmarshalJsonReaderMultiArray(t *testing.T) {
	var res struct {
		A string     `json:"a"`
		B [][]string `json:"b"`
	}
	payload := `{"a": "133", "b": [["add", "cccd"], ["eeee"]]}`
	reader := strings.NewReader(payload)
	if assert.NoError(t, UnmarshalJsonReader(reader, &res)) {
		assert.Equal(t, 2, len(res.B))
	}
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
		if assert.NoError(t, UnmarshalJsonMap(map[string]interface{}{
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

func TestUnmarshalJsonReaderWithTypeMismatchString(t *testing.T) {
	var req struct {
		Params map[string]string `json:"params"`
	}
	body := `{"params":{"a":{"a":123}}}`
	assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
}

func TestUnmarshalJsonReaderWithMismatchType(t *testing.T) {
	type Req struct {
		Params map[string]string `json:"params"`
	}

	var req Req
	body := `{"params":{"a":{"a":123}}}`
	assert.Equal(t, errTypeMismatch, UnmarshalJsonReader(strings.NewReader(body), &req))
}

func TestUnmarshalJsonReaderWithMismatchTypeBool(t *testing.T) {
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
}

func TestUnmarshalJsonReaderWithMismatchTypeBoolMap(t *testing.T) {
	var req struct {
		Params map[string]string `json:"params"`
	}
	assert.Equal(t, errTypeMismatch, UnmarshalJsonMap(map[string]interface{}{
		"params": map[string]interface{}{
			"a": true,
		},
	}, &req))
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
	m := map[string]interface{}{
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
	m := map[string]interface{}{
		"int": 1,
	}

	in := new(inner)
	if assert.NoError(t, UnmarshalKey(m, &in)) {
		assert.Equal(t, 1, in.Int)
	}
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
	m := map[string]interface{}{
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

	m := map[string]interface{}{
		"Ids": []map[string]interface{}{
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
	data := map[string]interface{}{
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
	data := map[string]interface{}{
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
