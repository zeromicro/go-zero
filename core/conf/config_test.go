package conf

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/fs"
	"github.com/zeromicro/go-zero/core/hash"
)

var dupErr conflictKeyError

func TestLoadConfig_notExists(t *testing.T) {
	assert.NotNil(t, Load("not_a_file", nil))
}

func TestLoadConfig_notRecogFile(t *testing.T) {
	filename, err := fs.TempFilenameWithText("hello")
	assert.Nil(t, err)
	defer os.Remove(filename)
	assert.NotNil(t, LoadConfig(filename, nil))
}

func TestConfigJson(t *testing.T) {
	tests := []string{
		".json",
		".yaml",
		".yml",
	}
	text := `{
	"a": "foo",
	"b": 1,
	"c": "${FOO}",
	"d": "abcd!@#$112"
}`
	t.Setenv("FOO", "2")

	for _, test := range tests {
		test := test
		t.Run(test, func(t *testing.T) {
			tmpfile, err := createTempFile(t, test, text)
			assert.Nil(t, err)

			var val struct {
				A string `json:"a"`
				B int    `json:"b"`
				C string `json:"c"`
				D string `json:"d"`
			}
			MustLoad(tmpfile, &val)
			assert.Equal(t, "foo", val.A)
			assert.Equal(t, 1, val.B)
			assert.Equal(t, "${FOO}", val.C)
			assert.Equal(t, "abcd!@#$112", val.D)
		})
	}
}

func TestLoadFromJsonBytesArray(t *testing.T) {
	input := []byte(`{"users": [{"name": "foo"}, {"Name": "bar"}]}`)
	var val struct {
		Users []struct {
			Name string
		}
	}

	assert.NoError(t, LoadConfigFromJsonBytes(input, &val))
	var expect []string
	for _, user := range val.Users {
		expect = append(expect, user.Name)
	}
	assert.EqualValues(t, []string{"foo", "bar"}, expect)
}

func TestConfigToml(t *testing.T) {
	text := `a = "foo"
b = 1
c = "${FOO}"
d = "abcd!@#$112"
`
	t.Setenv("FOO", "2")
	tmpfile, err := createTempFile(t, ".toml", text)
	assert.Nil(t, err)

	var val struct {
		A string `json:"a"`
		B int    `json:"b"`
		C string `json:"c"`
		D string `json:"d"`
	}
	MustLoad(tmpfile, &val)
	assert.Equal(t, "foo", val.A)
	assert.Equal(t, 1, val.B)
	assert.Equal(t, "${FOO}", val.C)
	assert.Equal(t, "abcd!@#$112", val.D)
}

func TestConfigOptional(t *testing.T) {
	text := `a = "foo"
b = 1
c = "FOO"
d = "abcd"
`
	tmpfile, err := createTempFile(t, ".toml", text)
	assert.Nil(t, err)

	var val struct {
		A string `json:"a"`
		B int    `json:"b,optional"`
		C string `json:"c,optional=B"`
		D string `json:"d,optional=b"`
	}
	if assert.NoError(t, Load(tmpfile, &val)) {
		assert.Equal(t, "foo", val.A)
		assert.Equal(t, 1, val.B)
		assert.Equal(t, "FOO", val.C)
		assert.Equal(t, "abcd", val.D)
	}
}

func TestConfigWithLower(t *testing.T) {
	text := `a = "foo"
b = 1
`
	tmpfile, err := createTempFile(t, ".toml", text)
	assert.Nil(t, err)

	var val struct {
		A string `json:"a"`
		b int
	}
	if assert.NoError(t, Load(tmpfile, &val)) {
		assert.Equal(t, "foo", val.A)
		assert.Equal(t, 0, val.b)
	}
}

func TestConfigJsonCanonical(t *testing.T) {
	text := []byte(`{"a": "foo", "B": "bar"}`)

	var val1 struct {
		A string `json:"a"`
		B string `json:"b"`
	}
	var val2 struct {
		A string
		B string
	}
	assert.NoError(t, LoadFromJsonBytes(text, &val1))
	assert.Equal(t, "foo", val1.A)
	assert.Equal(t, "bar", val1.B)
	assert.NoError(t, LoadFromJsonBytes(text, &val2))
	assert.Equal(t, "foo", val2.A)
	assert.Equal(t, "bar", val2.B)
}

func TestConfigTomlCanonical(t *testing.T) {
	text := []byte(`a = "foo"
B = "bar"`)

	var val1 struct {
		A string `json:"a"`
		B string `json:"b"`
	}
	var val2 struct {
		A string
		B string
	}
	assert.NoError(t, LoadFromTomlBytes(text, &val1))
	assert.Equal(t, "foo", val1.A)
	assert.Equal(t, "bar", val1.B)
	assert.NoError(t, LoadFromTomlBytes(text, &val2))
	assert.Equal(t, "foo", val2.A)
	assert.Equal(t, "bar", val2.B)
}

func TestConfigYamlCanonical(t *testing.T) {
	text := []byte(`a: foo
B: bar`)

	var val1 struct {
		A string `json:"a"`
		B string `json:"b"`
	}
	var val2 struct {
		A string
		B string
	}
	assert.NoError(t, LoadConfigFromYamlBytes(text, &val1))
	assert.Equal(t, "foo", val1.A)
	assert.Equal(t, "bar", val1.B)
	assert.NoError(t, LoadFromYamlBytes(text, &val2))
	assert.Equal(t, "foo", val2.A)
	assert.Equal(t, "bar", val2.B)
}

func TestConfigTomlEnv(t *testing.T) {
	text := `a = "foo"
b = 1
c = "${FOO}"
d = "abcd!@#112"
`
	t.Setenv("FOO", "2")
	tmpfile, err := createTempFile(t, ".toml", text)
	assert.Nil(t, err)

	var val struct {
		A string `json:"a"`
		B int    `json:"b"`
		C string `json:"c"`
		D string `json:"d"`
	}

	MustLoad(tmpfile, &val, UseEnv())
	assert.Equal(t, "foo", val.A)
	assert.Equal(t, 1, val.B)
	assert.Equal(t, "2", val.C)
	assert.Equal(t, "abcd!@#112", val.D)
}

func TestConfigJsonEnv(t *testing.T) {
	tests := []string{
		".json",
		".yaml",
		".yml",
	}
	text := `{
	"a": "foo",
	"b": 1,
	"c": "${FOO}",
	"d": "abcd!@#$a12 3"
}`
	t.Setenv("FOO", "2")
	for _, test := range tests {
		test := test
		t.Run(test, func(t *testing.T) {
			tmpfile, err := createTempFile(t, test, text)
			assert.Nil(t, err)

			var val struct {
				A string `json:"a"`
				B int    `json:"b"`
				C string `json:"c"`
				D string `json:"d"`
			}
			MustLoad(tmpfile, &val, UseEnv())
			assert.Equal(t, "foo", val.A)
			assert.Equal(t, 1, val.B)
			assert.Equal(t, "2", val.C)
			assert.Equal(t, "abcd!@# 3", val.D)
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  "",
			expect: "",
		},
		{
			input:  "A",
			expect: "a",
		},
		{
			input:  "a",
			expect: "a",
		},
		{
			input:  "hello_world",
			expect: "hello_world",
		},
		{
			input:  "Hello_world",
			expect: "hello_world",
		},
		{
			input:  "hello_World",
			expect: "hello_world",
		},
		{
			input:  "helloWorld",
			expect: "helloworld",
		},
		{
			input:  "HelloWorld",
			expect: "helloworld",
		},
		{
			input:  "hello World",
			expect: "hello world",
		},
		{
			input:  "Hello World",
			expect: "hello world",
		},
		{
			input:  "Hello World",
			expect: "hello world",
		},
		{
			input:  "Hello World foo_bar",
			expect: "hello world foo_bar",
		},
		{
			input:  "Hello World foo_Bar",
			expect: "hello world foo_bar",
		},
		{
			input:  "Hello World Foo_bar",
			expect: "hello world foo_bar",
		},
		{
			input:  "Hello World Foo_Bar",
			expect: "hello world foo_bar",
		},
		{
			input:  "Hello.World Foo_Bar",
			expect: "hello.world foo_bar",
		},
		{
			input:  "你好 World Foo_Bar",
			expect: "你好 world foo_bar",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			assert.Equal(t, test.expect, toLowerCase(test.input))
		})
	}
}

func TestLoadFromJsonBytesError(t *testing.T) {
	var val struct{}
	assert.Error(t, LoadFromJsonBytes([]byte(`hello`), &val))
}

func TestLoadFromTomlBytesError(t *testing.T) {
	var val struct{}
	assert.Error(t, LoadFromTomlBytes([]byte(`hello`), &val))
}

func TestLoadFromYamlBytesError(t *testing.T) {
	var val struct{}
	assert.Error(t, LoadFromYamlBytes([]byte(`':hello`), &val))
}

func TestLoadFromYamlBytes(t *testing.T) {
	input := []byte(`layer1:
  layer2:
    layer3: foo`)
	var val struct {
		Layer1 struct {
			Layer2 struct {
				Layer3 string
			}
		}
	}

	assert.NoError(t, LoadFromYamlBytes(input, &val))
	assert.Equal(t, "foo", val.Layer1.Layer2.Layer3)
}

func TestLoadFromYamlBytesTerm(t *testing.T) {
	input := []byte(`layer1:
  layer2:
    tls_conf: foo`)
	var val struct {
		Layer1 struct {
			Layer2 struct {
				Layer3 string `json:"tls_conf"`
			}
		}
	}

	assert.NoError(t, LoadFromYamlBytes(input, &val))
	assert.Equal(t, "foo", val.Layer1.Layer2.Layer3)
}

func TestLoadFromYamlBytesLayers(t *testing.T) {
	input := []byte(`layer1:
  layer2:
    layer3: foo`)
	var val struct {
		Value string `json:"Layer1.Layer2.Layer3"`
	}

	assert.NoError(t, LoadFromYamlBytes(input, &val))
	assert.Equal(t, "foo", val.Value)
}

func TestLoadFromYamlItemOverlay(t *testing.T) {
	type (
		Redis struct {
			Host string
			Port int
		}

		RedisKey struct {
			Redis
			Key string
		}

		Server struct {
			Redis RedisKey
		}

		TestConfig struct {
			Server
			Redis Redis
		}
	)

	input := []byte(`Redis:
  Host: localhost
  Port: 6379
  Key: test
`)

	var c TestConfig
	assert.ErrorAs(t, LoadFromYamlBytes(input, &c), &dupErr)
}

func TestLoadFromYamlItemOverlayReverse(t *testing.T) {
	type (
		Redis struct {
			Host string
			Port int
		}

		RedisKey struct {
			Redis
			Key string
		}

		Server struct {
			Redis Redis
		}

		TestConfig struct {
			Redis RedisKey
			Server
		}
	)

	input := []byte(`Redis:
  Host: localhost
  Port: 6379
  Key: test
`)

	var c TestConfig
	assert.ErrorAs(t, LoadFromYamlBytes(input, &c), &dupErr)
}

func TestLoadFromYamlItemOverlayWithMap(t *testing.T) {
	type (
		Redis struct {
			Host string
			Port int
		}

		RedisKey struct {
			Redis
			Key string
		}

		Server struct {
			Redis RedisKey
		}

		TestConfig struct {
			Server
			Redis map[string]interface{}
		}
	)

	input := []byte(`Redis:
  Host: localhost
  Port: 6379
  Key: test
`)

	var c TestConfig
	assert.ErrorAs(t, LoadFromYamlBytes(input, &c), &dupErr)
}

func TestUnmarshalJsonBytesMap(t *testing.T) {
	input := []byte(`{"foo":{"/mtproto.RPCTos": "bff.bff","bar":"baz"}}`)

	var val struct {
		Foo map[string]string
	}

	assert.NoError(t, LoadFromJsonBytes(input, &val))
	assert.Equal(t, "bff.bff", val.Foo["/mtproto.RPCTos"])
	assert.Equal(t, "baz", val.Foo["bar"])
}

func TestUnmarshalJsonBytesMapWithSliceElements(t *testing.T) {
	input := []byte(`{"foo":{"/mtproto.RPCTos": ["bff.bff", "any"],"bar":["baz", "qux"]}}`)

	var val struct {
		Foo map[string][]string
	}

	assert.NoError(t, LoadFromJsonBytes(input, &val))
	assert.EqualValues(t, []string{"bff.bff", "any"}, val.Foo["/mtproto.RPCTos"])
	assert.EqualValues(t, []string{"baz", "qux"}, val.Foo["bar"])
}

func TestUnmarshalJsonBytesMapWithSliceOfStructs(t *testing.T) {
	input := []byte(`{"foo":{
	"/mtproto.RPCTos": [{"bar": "any"}],
	"bar":[{"bar": "qux"}, {"bar": "ever"}]}}`)

	var val struct {
		Foo map[string][]struct {
			Bar string
		}
	}

	assert.NoError(t, LoadFromJsonBytes(input, &val))
	assert.Equal(t, 1, len(val.Foo["/mtproto.RPCTos"]))
	assert.Equal(t, "any", val.Foo["/mtproto.RPCTos"][0].Bar)
	assert.Equal(t, 2, len(val.Foo["bar"]))
	assert.Equal(t, "qux", val.Foo["bar"][0].Bar)
	assert.Equal(t, "ever", val.Foo["bar"][1].Bar)
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
		input = []byte(`{"Name": "hello", "int": 3}`)
		c     Conf
	)
	assert.NoError(t, LoadFromJsonBytes(input, &c))
	assert.Equal(t, "hello", c.Name)
	assert.Equal(t, Int(3), c.Int)
}

func TestUnmarshalJsonBytesWithMapValueOfStruct(t *testing.T) {
	type (
		Value struct {
			Name string
		}

		Config struct {
			Items map[string]Value
		}
	)

	var inputs = [][]byte{
		[]byte(`{"Items": {"Key":{"Name": "foo"}}}`),
		[]byte(`{"Items": {"Key":{"Name": "foo"}}}`),
		[]byte(`{"items": {"key":{"name": "foo"}}}`),
		[]byte(`{"items": {"key":{"name": "foo"}}}`),
	}
	for _, input := range inputs {
		var c Config
		if assert.NoError(t, LoadFromJsonBytes(input, &c)) {
			assert.Equal(t, 1, len(c.Items))
			for _, v := range c.Items {
				assert.Equal(t, "foo", v.Name)
			}
		}
	}
}

func TestUnmarshalJsonBytesWithMapTypeValueOfStruct(t *testing.T) {
	type (
		Value struct {
			Name string
		}

		Map map[string]Value

		Config struct {
			Map
		}
	)

	var inputs = [][]byte{
		[]byte(`{"Map": {"Key":{"Name": "foo"}}}`),
		[]byte(`{"Map": {"Key":{"Name": "foo"}}}`),
		[]byte(`{"map": {"key":{"name": "foo"}}}`),
		[]byte(`{"map": {"key":{"name": "foo"}}}`),
	}
	for _, input := range inputs {
		var c Config
		if assert.NoError(t, LoadFromJsonBytes(input, &c)) {
			assert.Equal(t, 1, len(c.Map))
			for _, v := range c.Map {
				assert.Equal(t, "foo", v.Name)
			}
		}
	}
}

func Test_FieldOverwrite(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		type Base struct {
			Name string
		}

		type St1 struct {
			Base
			Name2 string
		}

		type St2 struct {
			Base
			Name2 string
		}

		type St3 struct {
			*Base
			Name2 string
		}

		type St4 struct {
			*Base
			Name2 *string
		}

		validate := func(val any) {
			input := []byte(`{"Name": "hello", "Name2": "world"}`)
			assert.NoError(t, LoadFromJsonBytes(input, val))
		}

		validate(&St1{})
		validate(&St2{})
		validate(&St3{})
		validate(&St4{})
	})

	t.Run("Inherit Override", func(t *testing.T) {
		type Base struct {
			Name string
		}

		type St1 struct {
			Base
			Name string
		}

		type St2 struct {
			Base
			Name int
		}

		type St3 struct {
			*Base
			Name int
		}

		type St4 struct {
			*Base
			Name *string
		}

		validate := func(val any) {
			input := []byte(`{"Name": "hello"}`)
			err := LoadFromJsonBytes(input, val)
			assert.ErrorAs(t, err, &dupErr)
			assert.Equal(t, newConflictKeyError("name").Error(), err.Error())
		}

		validate(&St1{})
		validate(&St2{})
		validate(&St3{})
		validate(&St4{})
	})

	t.Run("Inherit more", func(t *testing.T) {
		type Base1 struct {
			Name string
		}

		type St0 struct {
			Base1
			Name string
		}

		type St1 struct {
			St0
			Name string
		}

		type St2 struct {
			St0
			Name int
		}

		type St3 struct {
			*St0
			Name int
		}

		type St4 struct {
			*St0
			Name *int
		}

		validate := func(val any) {
			input := []byte(`{"Name": "hello"}`)
			err := LoadFromJsonBytes(input, val)
			assert.ErrorAs(t, err, &dupErr)
			assert.Error(t, err)
		}

		validate(&St0{})
		validate(&St1{})
		validate(&St2{})
		validate(&St3{})
		validate(&St4{})
	})
}

func TestFieldOverwriteComplicated(t *testing.T) {
	t.Run("double maps", func(t *testing.T) {
		type (
			Base1 struct {
				Values map[string]string
			}
			Base2 struct {
				Values map[string]string
			}
			Config struct {
				Base1
				Base2
			}
		)

		var c Config
		input := []byte(`{"Values": {"Key": "Value"}}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("merge children", func(t *testing.T) {
		type (
			Inner1 struct {
				Name string
			}
			Inner2 struct {
				Age int
			}
			Base1 struct {
				Inner Inner1
			}
			Base2 struct {
				Inner Inner2
			}
			Config struct {
				Base1
				Base2
			}
		)

		var c Config
		input := []byte(`{"Inner": {"Name": "foo", "Age": 10}}`)
		if assert.NoError(t, LoadFromJsonBytes(input, &c)) {
			assert.Equal(t, "foo", c.Base1.Inner.Name)
			assert.Equal(t, 10, c.Base2.Inner.Age)
		}
	})

	t.Run("overwritten maps", func(t *testing.T) {
		type (
			Inner struct {
				Map map[string]string
			}
			Config struct {
				Map map[string]string
				Inner
			}
		)

		var c Config
		input := []byte(`{"Inner": {"Map": {"Key": "Value"}}}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("overwritten nested maps", func(t *testing.T) {
		type (
			Inner struct {
				Map map[string]string
			}
			Middle1 struct {
				Map map[string]string
				Inner
			}
			Middle2 struct {
				Map map[string]string
				Inner
			}
			Config struct {
				Middle1
				Middle2
			}
		)

		var c Config
		input := []byte(`{"Middle1": {"Inner": {"Map": {"Key": "Value"}}}}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("overwritten outer/inner maps", func(t *testing.T) {
		type (
			Inner struct {
				Map map[string]string
			}
			Middle struct {
				Inner
				Map map[string]string
			}
			Config struct {
				Middle
			}
		)

		var c Config
		input := []byte(`{"Middle": {"Inner": {"Map": {"Key": "Value"}}}}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("overwritten anonymous maps", func(t *testing.T) {
		type (
			Inner struct {
				Map map[string]string
			}
			Middle struct {
				Inner
				Map map[string]string
			}
			Elem   map[string]Middle
			Config struct {
				Elem
			}
		)

		var c Config
		input := []byte(`{"Elem": {"Key": {"Inner": {"Map": {"Key": "Value"}}}}}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("overwritten primitive and map", func(t *testing.T) {
		type (
			Inner struct {
				Value string
			}
			Elem  map[string]Inner
			Named struct {
				Elem string
			}
			Config struct {
				Named
				Elem
			}
		)

		var c Config
		input := []byte(`{"Elem": {"Key": {"Value": "Value"}}}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("overwritten map and slice", func(t *testing.T) {
		type (
			Inner struct {
				Value string
			}
			Elem  []Inner
			Named struct {
				Elem string
			}
			Config struct {
				Named
				Elem
			}
		)

		var c Config
		input := []byte(`{"Elem": {"Key": {"Value": "Value"}}}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("overwritten map and string", func(t *testing.T) {
		type (
			Elem  string
			Named struct {
				Elem string
			}
			Config struct {
				Named
				Elem
			}
		)

		var c Config
		input := []byte(`{"Elem": {"Key": {"Value": "Value"}}}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})
}

func TestLoadNamedFieldOverwritten(t *testing.T) {
	t.Run("overwritten named struct", func(t *testing.T) {
		type (
			Elem  string
			Named struct {
				Elem string
			}
			Base struct {
				Named
				Elem
			}
			Config struct {
				Val Base
			}
		)

		var c Config
		input := []byte(`{"Val": {"Elem": {"Key": {"Value": "Value"}}}}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("overwritten named []struct", func(t *testing.T) {
		type (
			Elem  string
			Named struct {
				Elem string
			}
			Base struct {
				Named
				Elem
			}
			Config struct {
				Vals []Base
			}
		)

		var c Config
		input := []byte(`{"Vals": [{"Elem": {"Key": {"Value": "Value"}}}]}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("overwritten named map[string]struct", func(t *testing.T) {
		type (
			Elem  string
			Named struct {
				Elem string
			}
			Base struct {
				Named
				Elem
			}
			Config struct {
				Vals map[string]Base
			}
		)

		var c Config
		input := []byte(`{"Vals": {"Key": {"Elem": {"Key": {"Value": "Value"}}}}}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("overwritten named *struct", func(t *testing.T) {
		type (
			Elem  string
			Named struct {
				Elem string
			}
			Base struct {
				Named
				Elem
			}
			Config struct {
				Vals *Base
			}
		)

		var c Config
		input := []byte(`{"Vals": [{"Elem": {"Key": {"Value": "Value"}}}]}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("overwritten named struct", func(t *testing.T) {
		type (
			Named struct {
				Elem string
			}
			Base struct {
				Named
				Elem Named
			}
			Config struct {
				Val Base
			}
		)

		var c Config
		input := []byte(`{"Val": {"Elem": "Value"}}`)
		assert.ErrorAs(t, LoadFromJsonBytes(input, &c), &dupErr)
	})

	t.Run("overwritten named struct", func(t *testing.T) {
		type Config struct {
			Val chan int
		}

		var c Config
		input := []byte(`{"Val": 1}`)
		assert.Error(t, LoadFromJsonBytes(input, &c))
	})
}

func TestLoadLowerMemberShouldNotConflict(t *testing.T) {
	type (
		Redis struct {
			db uint
		}

		Config struct {
			db uint
			Redis
		}
	)

	var c Config
	assert.NoError(t, LoadFromJsonBytes([]byte(`{}`), &c))
	assert.Zero(t, c.db)
	assert.Zero(t, c.Redis.db)
}

func TestFillDefaultUnmarshal(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		type St struct{}
		err := FillDefault(St{})
		assert.Error(t, err)
	})

	t.Run("not nil", func(t *testing.T) {
		type St struct{}
		err := FillDefault(&St{})
		assert.NoError(t, err)
	})

	t.Run("default", func(t *testing.T) {
		type St struct {
			A string `json:",default=a"`
			B string
		}
		var st St
		err := FillDefault(&st)
		assert.NoError(t, err)
		assert.Equal(t, st.A, "a")
	})

	t.Run("env", func(t *testing.T) {
		type St struct {
			A string `json:",default=a"`
			B string
			C string `json:",env=TEST_C"`
		}
		t.Setenv("TEST_C", "c")

		var st St
		err := FillDefault(&st)
		assert.NoError(t, err)
		assert.Equal(t, st.A, "a")
		assert.Equal(t, st.C, "c")
	})

	t.Run("has value", func(t *testing.T) {
		type St struct {
			A string `json:",default=a"`
			B string
		}
		var st = St{
			A: "b",
		}
		err := FillDefault(&st)
		assert.Error(t, err)
	})
}

func TestConfigWithJsonTag(t *testing.T) {
	t.Run("map with value", func(t *testing.T) {
		var input = []byte(`[Value]
[Value.first]
Email = "foo"
[Value.second]
Email = "bar"`)

		type Value struct {
			Email string
		}

		type Config struct {
			ValueMap map[string]Value `json:"Value"`
		}

		var c Config
		if assert.NoError(t, LoadFromTomlBytes(input, &c)) {
			assert.Len(t, c.ValueMap, 2)
		}
	})

	t.Run("map with ptr value", func(t *testing.T) {
		var input = []byte(`[Value]
[Value.first]
Email = "foo"
[Value.second]
Email = "bar"`)

		type Value struct {
			Email string
		}

		type Config struct {
			ValueMap map[string]*Value `json:"Value"`
		}

		var c Config
		if assert.NoError(t, LoadFromTomlBytes(input, &c)) {
			assert.Len(t, c.ValueMap, 2)
		}
	})

	t.Run("map with optional", func(t *testing.T) {
		var input = []byte(`[Value]
[Value.first]
Email = "foo"
[Value.second]
Email = "bar"`)

		type Value struct {
			Email string
		}

		type Config struct {
			Value map[string]Value `json:",optional"`
		}

		var c Config
		if assert.NoError(t, LoadFromTomlBytes(input, &c)) {
			assert.Len(t, c.Value, 2)
		}
	})

	t.Run("map with empty tag", func(t *testing.T) {
		var input = []byte(`[Value]
[Value.first]
Email = "foo"
[Value.second]
Email = "bar"`)

		type Value struct {
			Email string
		}

		type Config struct {
			Value map[string]Value `json:"  "`
		}

		var c Config
		if assert.NoError(t, LoadFromTomlBytes(input, &c)) {
			assert.Len(t, c.Value, 2)
		}
	})

	t.Run("multi layer map", func(t *testing.T) {
		type Value struct {
			User struct {
				Name string
			}
		}

		type Config struct {
			Value map[string]map[string]Value
		}

		var input = []byte(`
[Value.first.User1.User]
Name = "foo"
[Value.second.User2.User]
Name = "bar"
`)
		var c Config
		if assert.NoError(t, LoadFromTomlBytes(input, &c)) {
			assert.Len(t, c.Value, 2)
		}
	})
}

func Test_LoadBadConfig(t *testing.T) {
	type Config struct {
		Name string `json:"name,options=foo|bar"`
	}

	file, err := createTempFile(t, ".json", `{"name": "baz"}`)
	assert.NoError(t, err)

	var c Config
	err = Load(file, &c)
	assert.Error(t, err)
}

func Test_getFullName(t *testing.T) {
	assert.Equal(t, "a.b", getFullName("a", "b"))
	assert.Equal(t, "a", getFullName("", "a"))
}

func TestValidate(t *testing.T) {
	t.Run("normal config", func(t *testing.T) {
		var c mockConfig
		err := LoadFromJsonBytes([]byte(`{"val": "hello", "number": 8}`), &c)
		assert.NoError(t, err)
	})

	t.Run("error no int", func(t *testing.T) {
		var c mockConfig
		err := LoadFromJsonBytes([]byte(`{"val": "hello"}`), &c)
		assert.Error(t, err)
	})

	t.Run("error no string", func(t *testing.T) {
		var c mockConfig
		err := LoadFromJsonBytes([]byte(`{"number": 8}`), &c)
		assert.Error(t, err)
	})
}

func Test_buildFieldsInfo(t *testing.T) {
	type ParentSt struct {
		Name string
		M    map[string]int
	}
	tests := []struct {
		name        string
		t           reflect.Type
		ok          bool
		containsKey string
	}{
		{
			name: "normal",
			t:    reflect.TypeOf(struct{ A string }{}),
			ok:   true,
		},
		{
			name: "struct anonymous",
			t: reflect.TypeOf(struct {
				ParentSt
				Name string
			}{}),
			ok:          false,
			containsKey: newConflictKeyError("name").Error(),
		},
		{
			name: "struct ptr anonymous",
			t: reflect.TypeOf(struct {
				*ParentSt
				Name string
			}{}),
			ok:          false,
			containsKey: newConflictKeyError("name").Error(),
		},
		{
			name: "more struct anonymous",
			t: reflect.TypeOf(struct {
				Value struct {
					ParentSt
					Name string
				}
			}{}),
			ok:          false,
			containsKey: newConflictKeyError("value.name").Error(),
		},
		{
			name: "map anonymous",
			t: reflect.TypeOf(struct {
				ParentSt
				M string
			}{}),
			ok:          false,
			containsKey: newConflictKeyError("m").Error(),
		},
		{
			name: "map more anonymous",
			t: reflect.TypeOf(struct {
				Value struct {
					ParentSt
					M string
				}
			}{}),
			ok:          false,
			containsKey: newConflictKeyError("value.m").Error(),
		},
		{
			name: "struct slice anonymous",
			t: reflect.TypeOf([]struct {
				ParentSt
				Name string
			}{}),
			ok:          false,
			containsKey: newConflictKeyError("name").Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildFieldsInfo(tt.t, "")
			if tt.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tt.containsKey)
			}
		})
	}
}

func createTempFile(t *testing.T, ext, text string) (string, error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), hash.Md5Hex([]byte(text))+"*"+ext)
	if err != nil {
		return "", err
	}

	if err = os.WriteFile(tmpFile.Name(), []byte(text), os.ModeTemporary); err != nil {
		return "", err
	}

	filename := tmpFile.Name()
	if err = tmpFile.Close(); err != nil {
		return "", err
	}

	t.Cleanup(func() {
		_ = os.Remove(filename)
	})

	return filename, nil
}

type mockConfig struct {
	Val    string
	Number int
}

func (m mockConfig) Validate() error {
	if len(m.Val) == 0 {
		return errors.New("val is empty")
	}

	if m.Number == 0 {
		return errors.New("number is zero")
	}

	return nil
}
