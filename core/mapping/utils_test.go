package mapping

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testTagName = "key"

type Foo struct {
	Str                 string
	StrWithTag          string `key:"stringwithtag"`
	StrWithTagAndOption string `key:"stringwithtag,string"`
}

func TestDerefInt(t *testing.T) {
	i := 1
	s := "hello"
	number := struct {
		f float64
	}{
		f: 6.4,
	}
	cases := []struct {
		t      reflect.Type
		expect reflect.Kind
	}{
		{
			t:      reflect.TypeOf(i),
			expect: reflect.Int,
		},
		{
			t:      reflect.TypeOf(&i),
			expect: reflect.Int,
		},
		{
			t:      reflect.TypeOf(s),
			expect: reflect.String,
		},
		{
			t:      reflect.TypeOf(&s),
			expect: reflect.String,
		},
		{
			t:      reflect.TypeOf(number.f),
			expect: reflect.Float64,
		},
		{
			t:      reflect.TypeOf(&number.f),
			expect: reflect.Float64,
		},
	}

	for _, each := range cases {
		t.Run(each.t.String(), func(t *testing.T) {
			assert.Equal(t, each.expect, Deref(each.t).Kind())
		})
	}
}

func TestDerefValInt(t *testing.T) {
	i := 1
	s := "hello"
	number := struct {
		f float64
	}{
		f: 6.4,
	}
	cases := []struct {
		t      reflect.Value
		expect reflect.Kind
	}{
		{
			t:      reflect.ValueOf(i),
			expect: reflect.Int,
		},
		{
			t:      reflect.ValueOf(&i),
			expect: reflect.Int,
		},
		{
			t:      reflect.ValueOf(s),
			expect: reflect.String,
		},
		{
			t:      reflect.ValueOf(&s),
			expect: reflect.String,
		},
		{
			t:      reflect.ValueOf(number.f),
			expect: reflect.Float64,
		},
		{
			t:      reflect.ValueOf(&number.f),
			expect: reflect.Float64,
		},
	}

	for _, each := range cases {
		t.Run(each.t.String(), func(t *testing.T) {
			assert.Equal(t, each.expect, ensureValue(each.t).Kind())
		})
	}
}

func TestParseKeyAndOptionWithoutTag(t *testing.T) {
	var foo Foo
	rte := reflect.TypeOf(&foo).Elem()
	field, _ := rte.FieldByName("Str")
	key, options, err := parseKeyAndOptions(testTagName, field)
	assert.Nil(t, err)
	assert.Equal(t, "Str", key)
	assert.Nil(t, options)
}

func TestParseKeyAndOptionWithTagWithoutOption(t *testing.T) {
	var foo Foo
	rte := reflect.TypeOf(&foo).Elem()
	field, _ := rte.FieldByName("StrWithTag")
	key, options, err := parseKeyAndOptions(testTagName, field)
	assert.Nil(t, err)
	assert.Equal(t, "stringwithtag", key)
	assert.Nil(t, options)
}

func TestParseKeyAndOptionWithTagAndOption(t *testing.T) {
	var foo Foo
	rte := reflect.TypeOf(&foo).Elem()
	field, _ := rte.FieldByName("StrWithTagAndOption")
	key, options, err := parseKeyAndOptions(testTagName, field)
	assert.Nil(t, err)
	assert.Equal(t, "stringwithtag", key)
	assert.True(t, options.FromString)
}

func TestParseSegments(t *testing.T) {
	tests := []struct {
		input  string
		expect []string
	}{
		{
			input:  "",
			expect: []string{},
		},
		{
			input:  "   ",
			expect: []string{},
		},
		{
			input:  ",",
			expect: []string{""},
		},
		{
			input:  "foo,",
			expect: []string{"foo"},
		},
		{
			input: ",foo",
			// the first empty string cannot be ignored, it's the key.
			expect: []string{"", "foo"},
		},
		{
			input:  "foo",
			expect: []string{"foo"},
		},
		{
			input:  "foo,bar",
			expect: []string{"foo", "bar"},
		},
		{
			input:  "foo,bar,baz",
			expect: []string{"foo", "bar", "baz"},
		},
		{
			input:  "foo,options=a|b",
			expect: []string{"foo", "options=a|b"},
		},
		{
			input:  "foo,bar,default=[baz,qux]",
			expect: []string{"foo", "bar", "default=[baz,qux]"},
		},
		{
			input:  "foo,bar,options=[baz,qux]",
			expect: []string{"foo", "bar", "options=[baz,qux]"},
		},
		{
			input:  `foo\,bar,options=[baz,qux]`,
			expect: []string{`foo,bar`, "options=[baz,qux]"},
		},
		{
			input:  `foo,bar,options=\[baz,qux]`,
			expect: []string{"foo", "bar", "options=[baz", "qux]"},
		},
		{
			input:  `foo,bar,options=[baz\,qux]`,
			expect: []string{"foo", "bar", `options=[baz\,qux]`},
		},
		{
			input:  `foo\,bar,options=[baz,qux],default=baz`,
			expect: []string{`foo,bar`, "options=[baz,qux]", "default=baz"},
		},
		{
			input:  `foo\,bar,options=[baz,qux, quux],default=[qux, baz]`,
			expect: []string{`foo,bar`, "options=[baz,qux, quux]", "default=[qux, baz]"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			assert.ElementsMatch(t, test.expect, parseSegments(test.input))
		})
	}
}

func TestValidatePtrWithNonPtr(t *testing.T) {
	var foo string
	rve := reflect.ValueOf(foo)
	assert.NotNil(t, ValidatePtr(rve))
}

func TestValidatePtrWithPtr(t *testing.T) {
	var foo string
	rve := reflect.ValueOf(&foo)
	assert.Nil(t, ValidatePtr(rve))
}

func TestValidatePtrWithNilPtr(t *testing.T) {
	var foo *string
	rve := reflect.ValueOf(foo)
	assert.NotNil(t, ValidatePtr(rve))
}

func TestValidatePtrWithZeroValue(t *testing.T) {
	var s string
	e := reflect.Zero(reflect.TypeOf(s))
	assert.NotNil(t, ValidatePtr(e))
}

func TestSetValueNotSettable(t *testing.T) {
	var i int
	assert.Error(t, setValueFromString(reflect.Int, reflect.ValueOf(i), "1"))
	assert.Error(t, validateAndSetValue(reflect.Int, reflect.ValueOf(i), "1", nil))
}

func TestParseKeyAndOptionsErrors(t *testing.T) {
	type Bar struct {
		OptionsValue string `key:",options=a=b"`
		DefaultValue string `key:",default=a=b"`
	}

	var bar Bar
	_, _, err := parseKeyAndOptions("key", reflect.TypeOf(&bar).Elem().Field(0))
	assert.NotNil(t, err)
	_, _, err = parseKeyAndOptions("key", reflect.TypeOf(&bar).Elem().Field(1))
	assert.NotNil(t, err)
}

func TestSetValueFormatErrors(t *testing.T) {
	type Bar struct {
		IntValue   int
		UintValue  uint
		FloatValue float32
		MapValue   map[string]any
	}

	var bar Bar
	tests := []struct {
		kind   reflect.Kind
		target reflect.Value
		value  string
	}{
		{
			kind:   reflect.Int,
			target: reflect.ValueOf(&bar.IntValue).Elem(),
			value:  "a",
		},
		{
			kind:   reflect.Uint,
			target: reflect.ValueOf(&bar.UintValue).Elem(),
			value:  "a",
		},
		{
			kind:   reflect.Float32,
			target: reflect.ValueOf(&bar.FloatValue).Elem(),
			value:  "a",
		},
		{
			kind:   reflect.Map,
			target: reflect.ValueOf(&bar.MapValue).Elem(),
		},
	}

	for _, test := range tests {
		t.Run(test.kind.String(), func(t *testing.T) {
			err := setValueFromString(test.kind, test.target, test.value)
			assert.NotEqual(t, errValueNotSettable, err)
			assert.NotNil(t, err)
		})
	}
}

func TestValidateValueRange(t *testing.T) {
	t.Run("float", func(t *testing.T) {
		assert.NoError(t, validateValueRange(1.2, nil))
	})

	t.Run("float number range", func(t *testing.T) {
		assert.NoError(t, validateNumberRange(1.2, nil))
	})

	t.Run("bad float", func(t *testing.T) {
		assert.Error(t, validateValueRange("a", &fieldOptionsWithContext{
			Range: &numberRange{},
		}))
	})

	t.Run("bad float validate", func(t *testing.T) {
		var v struct {
			Foo float32
		}
		assert.Error(t, validateAndSetValue(reflect.Int, reflect.ValueOf(&v).Elem().Field(0),
			"1", &fieldOptionsWithContext{
				Range: &numberRange{
					left:  2,
					right: 3,
				},
			}))
	})
}

func TestSetMatchedPrimitiveValue(t *testing.T) {
	assert.Error(t, setMatchedPrimitiveValue(reflect.Func, reflect.ValueOf(2), "1"))
}
