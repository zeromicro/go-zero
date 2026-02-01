package conf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/fs"
)

func TestProperties(t *testing.T) {
	text := `app.name = test

    app.program=app

    # this is comment
    app.threads = 5`
	tmpfile, err := fs.TempFilenameWithText(text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

	props, err := LoadProperties(tmpfile)
	assert.Nil(t, err)
	assert.Equal(t, "test", props.GetString("app.name"))
	assert.Equal(t, "app", props.GetString("app.program"))
	assert.Equal(t, 5, props.GetInt("app.threads"))

	val := props.ToString()
	assert.Contains(t, val, "app.name")
	assert.Contains(t, val, "app.program")
	assert.Contains(t, val, "app.threads")
}

func TestPropertiesEnv(t *testing.T) {
	text := `app.name = test

    app.program=app

	app.env1 = ${FOO}
	app.env2 = $none

    # this is comment
    app.threads = 5`
	tmpfile, err := fs.TempFilenameWithText(text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

	t.Setenv("FOO", "2")

	props, err := LoadProperties(tmpfile, UseEnv())
	assert.Nil(t, err)
	assert.Equal(t, "test", props.GetString("app.name"))
	assert.Equal(t, "app", props.GetString("app.program"))
	assert.Equal(t, 5, props.GetInt("app.threads"))
	assert.Equal(t, "2", props.GetString("app.env1"))
	assert.Equal(t, "", props.GetString("app.env2"))

	val := props.ToString()
	assert.Contains(t, val, "app.name")
	assert.Contains(t, val, "app.program")
	assert.Contains(t, val, "app.threads")
	assert.Contains(t, val, "app.env1")
	assert.Contains(t, val, "app.env2")
}

func TestLoadProperties_badContent(t *testing.T) {
	filename, err := fs.TempFilenameWithText("hello")
	assert.Nil(t, err)
	defer os.Remove(filename)
	_, err = LoadProperties(filename)
	assert.NotNil(t, err)
	assert.True(t, len(err.Error()) > 0)
}

func TestSetString(t *testing.T) {
	key := "a"
	value := "the value of a"
	props := NewProperties()
	props.SetString(key, value)
	assert.Equal(t, value, props.GetString(key))
}

func TestSetInt(t *testing.T) {
	key := "a"
	value := 101
	props := NewProperties()
	props.SetInt(key, value)
	assert.Equal(t, value, props.GetInt(key))
}

func TestLoadBadFile(t *testing.T) {
	_, err := LoadProperties("nosuchfile")
	assert.NotNil(t, err)
}

func TestProperties_valueWithEqualSymbols(t *testing.T) {
	text := `# test with equal symbols in value
	db.url=postgres://localhost:5432/db?param=value
	math.equation=a=b=c
	base64.data=SGVsbG8=World=Test=
	url.with.params=http://example.com?foo=bar&baz=qux
	empty.value=
	key.with.space = value = with = equals`
	tmpfile, err := fs.TempFilenameWithText(text)
	assert.Nil(t, err)
	defer os.Remove(tmpfile)

	props, err := LoadProperties(tmpfile)
	assert.Nil(t, err)
	assert.Equal(t, "postgres://localhost:5432/db?param=value", props.GetString("db.url"))
	assert.Equal(t, "a=b=c", props.GetString("math.equation"))
	assert.Equal(t, "SGVsbG8=World=Test=", props.GetString("base64.data"))
	assert.Equal(t, "http://example.com?foo=bar&baz=qux", props.GetString("url.with.params"))
	assert.Equal(t, "", props.GetString("empty.value"))
	assert.Equal(t, "value = with = equals", props.GetString("key.with.space"))
}

func TestProperties_edgeCases(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "no equal sign",
			content: "invalid line without equal",
			wantErr: true,
		},
		{
			name:    "only equal sign",
			content: "=",
			wantErr: false, // "=" 会被解析为空 key 和空 value，len(pair) == 2，是合法的
		},
		{
			name:    "empty key",
			content: "=value",
			wantErr: false, // 空 key 也会被 trim，但 len(pair) == 2 所以不会报错
		},
		{
			name:    "equal at end",
			content: "key.name=",
			wantErr: false, // 空 value 是合法的
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := fs.TempFilenameWithText(tt.content)
			assert.Nil(t, err)
			defer os.Remove(tmpfile)

			_, err = LoadProperties(tmpfile)
			if tt.wantErr {
				assert.NotNil(t, err, "expected error for case: %s", tt.name)
			} else {
				assert.Nil(t, err, "unexpected error for case: %s", tt.name)
			}
		})
	}
}
