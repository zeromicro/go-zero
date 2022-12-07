package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/fs"
	"github.com/zeromicro/go-zero/core/hash"
)

func TestLoadConfig_notExists(t *testing.T) {
	assert.NotNil(t, Load("not_a_file", nil))
}

func TestLoadConfig_notRecogFile(t *testing.T) {
	filename, err := fs.TempFilenameWithText("hello")
	assert.Nil(t, err)
	defer os.Remove(filename)
	assert.NotNil(t, Load(filename, nil))
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
	for _, test := range tests {
		test := test
		t.Run(test, func(t *testing.T) {
			os.Setenv("FOO", "2")
			defer os.Unsetenv("FOO")
			tmpfile, err := createTempFile(test, text)
			assert.Nil(t, err)
			defer os.Remove(tmpfile)

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

	assert.NoError(t, LoadFromJsonBytes(input, &val))
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
	os.Setenv("FOO", "2")
	defer os.Unsetenv("FOO")
	tmpfile, err := createTempFile(".toml", text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

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
	assert.NoError(t, LoadFromYamlBytes(text, &val1))
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
	os.Setenv("FOO", "2")
	defer os.Unsetenv("FOO")
	tmpfile, err := createTempFile(".toml", text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

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
	for _, test := range tests {
		test := test
		t.Run(test, func(t *testing.T) {
			os.Setenv("FOO", "2")
			defer os.Unsetenv("FOO")
			tmpfile, err := createTempFile(test, text)
			assert.Nil(t, err)
			defer os.Remove(tmpfile)

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
			expect: "helloWorld",
		},
		{
			input:  "Hello_world",
			expect: "helloWorld",
		},
		{
			input:  "hello_World",
			expect: "helloWorld",
		},
		{
			input:  "helloWorld",
			expect: "helloWorld",
		},
		{
			input:  "HelloWorld",
			expect: "helloWorld",
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
			expect: "hello world fooBar",
		},
		{
			input:  "Hello World foo_Bar",
			expect: "hello world fooBar",
		},
		{
			input:  "Hello World Foo_bar",
			expect: "hello world fooBar",
		},
		{
			input:  "Hello World Foo_Bar",
			expect: "hello world fooBar",
		},
		{
			input:  "你好 World Foo_Bar",
			expect: "你好 world fooBar",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.input, func(t *testing.T) {
			assert.Equal(t, test.expect, toCamelCase(test.input))
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

func createTempFile(ext, text string) (string, error) {
	tmpfile, err := os.CreateTemp(os.TempDir(), hash.Md5Hex([]byte(text))+"*"+ext)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(tmpfile.Name(), []byte(text), os.ModeTemporary); err != nil {
		return "", err
	}

	filename := tmpfile.Name()
	if err = tmpfile.Close(); err != nil {
		return "", err
	}

	return filename, nil
}
